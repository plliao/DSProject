package backEnd

import (
    "net"
    "net/http"
    "net/rpc"
    "log"
    "regexp"
    "errors"
    //"crypto/rand"
    "fmt"
    "reflect"
    "backEnd/cmd"
    "sync"
    "time"
    "strings"
    "bufio"
    "os"
)

type Server struct {
    users map[string]*User
    tokens map[string]*User
    commands chan reflect.Value
    cmdFactory *cmd.CommandFactory
    commandLogs map[string]reflect.Value

    rwLock *sync.RWMutex
    service *Service
    messages BackEndMessages

    validUserName *regexp.Regexp
    validPassword *regexp.Regexp

    network string
    addressBook []string
    raft *Raft
    commitChan chan int
    nextIndexs []int
    id int
    timeout time.Duration

    toExecChan chan int
    heartBeatChan chan time.Time
    lastBeatTime time.Time

    logger *bufio.Writer
}

func (srv *Server) Init(id int) {
    srv.users = make(map[string]*User)
    srv.tokens = make(map[string]*User)
    srv.commands = make(chan reflect.Value, 100)
    srv.cmdFactory = &cmd.CommandFactory{}
    srv.cmdFactory.Init()

    srv.rwLock = &sync.RWMutex{}
    srv.service = &Service{srv.commands}
    srv.messages = BackEndMessages{
        NoError:"",
        IncorrectFormat:"Incorrect username or password format",
        IncorrectPassword:"Incorrect password",
        UserNotExist:"User not exists",
        UserAlreadyExist:"User already exists",
        UnrecognizedToken:"Unrecognized token",
        FunctionNotImplement:"Function not implemented",
        EmptyToken:"",
        NotLeader:"Not Leader: ",
    }

    srv.validUserName, _ = regexp.Compile("^[a-zA-Z0-9]{4,10}$")
    srv.validPassword, _ = regexp.Compile("^[a-zA-Z0-9]{4,10}$")

    srv.id = id
    srv.network = "tcp"
    srv.addressBook = make([]string, 0)
    srv.timeout = 1000 * time.Millisecond
    srv.heartBeatChan = make(chan time.Time, 100)
}

func (srv *Server) RegisterAddress(address string, port string) {
    srv.addressBook = append(srv.addressBook, address + ":" + port)
    srv.nextIndexs = append(srv.nextIndexs, 0)
}

func (srv *Server) getMajority() int {
    return len(srv.addressBook) / 2
}

func (srv *Server) validateUserNameAndPassFormat(username string, password string) (bool, error) {
    if !srv.validUserName.Match([]byte(username)) || !srv.validPassword.Match([]byte(password)) {
        return false, errors.New(srv.messages.IncorrectFormat)
    }
    return true, nil
}

func (srv *Server) validateUser(username string, password string) (bool, error) {
    if ok, err := srv.validateUserNameAndPassFormat(username, password); !ok {
        return ok, err
    }
    if user, ok := srv.users[username]; ok {
        if user.Password != password {
            return false, errors.New(srv.messages.IncorrectPassword)
        }
        return true, nil
    }
    return false, errors.New(srv.messages.UserNotExist)
}

func (srv *Server) validateAuth(token string) bool {
    if _, ok := srv.tokens[token]; ok {
        return true
    }
    return false
}

func (srv *Server) generateUserToken(user *User) {
    srv.deleteUserToken(user)
    //token := make([]byte, 6)
    //rand.Read(token)
    //user.token = fmt.Sprintf("%x", token)
    user.token = user.Username + user.Password
    srv.tokens[user.token] = user
}

func (srv *Server) deleteUserToken(user *User) {
    if user.token != "" {
        delete(srv.tokens, user.token)
        user.token = ""
    }
}

func (srv *Server) removeUser(user *User) {
    srv.deleteUserToken(user)
    for _, follower := range user.followers {
        follower.UnFollow(user)
    }
    delete(srv.users, user.Username)
}

func (srv *Server) isReadOnly(funcName string) bool {
    switch (funcName) {
        case "GetFollower", "GetMyContent":
            return true
        default:
            return false
    }
}

func (srv *Server) replyWithResults(cmdValue reflect.Value, results []reflect.Value) {
    replyType := cmdValue.Field(1).Type().Elem().Elem()
    reply := reflect.New(replyType)

    for index, value := range results {
        reply.Elem().Field(index).Set(value)
    }
    fmt.Printf("\nReply Command: %v, reply: %v\n", cmdValue.Type().Name(), reply)
    if !cmdValue.Field(1).IsNil() {
        cmdValue.Field(1).Send(reply)
    }
}

func (srv *Server) replyNotLeader(cmdValue reflect.Value) {
    leaderAddress := ""
    if srv.raft.leaderId >= 0 {
        leaderAddress = srv.addressBook[srv.raft.leaderId]
    }
    results := make([]reflect.Value, 2)
    results[0] = reflect.ValueOf(false)
    results[1] = reflect.ValueOf(srv.messages.NotLeader + leaderAddress)
    srv.replyWithResults(cmdValue, results)
}

func (srv *Server) exec(encodedCmd string) []reflect.Value {
    srvValue := reflect.ValueOf(srv)
    funcName, parameters := srv.cmdFactory.Decode(encodedCmd)
    f := srvValue.MethodByName(funcName)

    if srv.isReadOnly(funcName) {
        srv.rwLock.RLock()
        defer srv.rwLock.RUnlock()
    } else {
        srv.rwLock.Lock()
        defer srv.rwLock.Unlock();
    }

    results := f.Call(parameters)
    return results
}

func (srv *Server) execAndReply(cmdValue reflect.Value) {
    encodedCmd := srv.cmdFactory.Encode(cmdValue)
    results := srv.exec(encodedCmd)
    srv.replyWithResults(cmdValue, results)
}

func (srv *Server) execCommit(commitIndex int) []reflect.Value {
    if srv.raft.index >= commitIndex {
        encodedCmd := srv.raft.logs[commitIndex]
        cmdTerm := srv.raft.logTerms[commitIndex]
        log := fmt.Sprintf("CommitIndex %v, term %v: %v\n", commitIndex, cmdTerm, encodedCmd)
        srv.logger.WriteString(log)
        srv.logger.Flush()
        return srv.exec(encodedCmd)
    }
    return nil
}

func (srv *Server) toFollowerHandler() {
    for term := range srv.raft.toFollowerChan {
        srv.rwLock.Lock()
        if srv.raft.term < term {
            srv.raft.term = term
            if srv.raft.isLeader {
                srv.leaderShutDown()
                srv.followerInit()
            }
        }
        srv.rwLock.Unlock()
    }
}

func (srv *Server) appendCommand(cmdValue reflect.Value) {
    encodedCmd := srv.cmdFactory.Encode(cmdValue)
    commandId := srv.cmdFactory.GetCommandId(encodedCmd)
    if _, ok := srv.raft.commandLogs[commandId]; !ok {
        srv.commandLogs[commandId] = cmdValue
        srv.raft.appendCommand(encodedCmd, srv.raft.term)
    } else if _, ok := srv.commandLogs[commandId]; !ok {
        srv.commandLogs[commandId] = cmdValue
        fmt.Printf("\nCommand %v exists\n", commandId)
    }
}

func (srv *Server) runCommands() {
    for cmd := range srv.commands {
        if !srv.raft.isLeader {
            srv.replyNotLeader(cmd)
            continue
        }
        if srv.isReadOnly(cmd.Type().Name()) {
            go srv.execAndReply(cmd)
        } else {
            srv.rwLock.RLock()
            if srv.raft.isLeader {
                srv.appendCommand(cmd)
            } else {
                srv.replyNotLeader(cmd)
            }
            srv.rwLock.RUnlock()
        }
    }
}

func (srv *Server) updateLastBeat() {
    for beatTime := range srv.raft.heartBeatChan {
        srv.lastBeatTime = beatTime
    }
}

func (srv *Server) Start() {
    go srv.runCommands()

    srv.raft = &Raft{
        isLeader:false,
        leaderId:-1,
        term:0,
        voteFor:-1,
        index:-1,
        commitIndex:-1,
        toFollowerChan:make(chan int, len(srv.addressBook)),
        heartBeatChan:srv.heartBeatChan,
        cmdFactory:srv.cmdFactory,
        commandLogs:make(map[string]int),
    }
    go srv.updateLastBeat()
    go srv.toFollowerHandler()
    address := srv.addressBook[srv.id]
    port := strings.Split(address, ":")[1]

    f, _ := os.Create("logs/" + address)
    defer f.Close()
    srv.logger = bufio.NewWriter(f)

    srv.followerInit()
    rpcServer := rpc.NewServer()
    rpcServer.Register(srv.service)
    rpcServer.Register(srv.raft)
    rpcServer.HandleHTTP(rpc.DefaultRPCPath + address, rpc.DefaultDebugPath + address)
    l, e := net.Listen("tcp", ":" + port)
    if e != nil {
        log.Fatal("listen error:", e)
    }
    fmt.Print("BackEnd serving on " + address + "\n")
    http.Serve(l, nil)
}

