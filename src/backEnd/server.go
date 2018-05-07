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
    srv.lastBeatTime = time.Now()
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
        case "GetMyFollower", "GetMyContent":
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
    fmt.Printf("Command: %v, reply: %v\n", cmdValue.Type().Name(), reply)
    cmdValue.Field(1).Send(reply)
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
    fmt.Printf("Execute command %v, %v\n", funcName, parameters)
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
        return srv.exec(encodedCmd)
    }
    return nil
}

func (srv *Server) toFollowerHandler() {
    for term := range srv.raft.toFollowerChan {
        if srv.raft.term >= term {
            continue
        }
        srv.rwLock.Lock()
        srv.raft.term = term
        if srv.raft.isLeader {
            srv.leaderShutDown()
            srv.followerInit()
        }
        srv.rwLock.Unlock()
    }
}

func (srv *Server) appendCommand(cmdValue reflect.Value) {
    encodedCmd := srv.cmdFactory.Encode(cmdValue)
    commandId := srv.cmdFactory.GetCommandId(encodedCmd)
    if _, ok := srv.commandLogs[commandId]; !ok {
        srv.commandLogs[commandId] = cmdValue
        srv.raft.appendCommand(encodedCmd, srv.raft.term)
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
            srv.appendCommand(cmd)
        }
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
    }
    go srv.toFollowerHandler()
    address := srv.addressBook[srv.id]
    port := strings.Split(address, ":")[1]

    srv.followerInit()
    rpc.Register(srv.service)
    rpc.Register(srv.raft)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":" + port)
    if e != nil {
        log.Fatal("listen error:", e)
    }
    fmt.Print("BackEnd serving on " + address + "\n")
    http.Serve(l, nil)
}

