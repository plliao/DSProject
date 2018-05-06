package backEnd

import (
    "net"
    "net/http"
    "net/rpc"
    "log"
    "regexp"
    "errors"
    //"crypto/rand"
    mrand "math/rand"
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
    commandLogs map[int]reflect.Value

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

func (srv *Server) leaderInit() {
    srv.raft.isLeader = true
    for i, _ := range srv.addressBook {
        srv.nextIndexs[i] = len(srv.raft.logs)
    }
    srv.commitChan = make(chan int, 100)
    srv.commandLogs = make(map[int]reflect.Value)
    go srv.commitHandler()
    for i, _ := range srv.addressBook {
        if i != srv.id {
            go srv.followerHandler(i)
        }
    }
}

func (srv *Server) leaderShutDown() {
    close(srv.commitChan)
    srv.commandLogs = nil
    srv.raft.isLeader = false
}

func (srv *Server) followerInit() {
    srv.toExecChan = make(chan int, 100)
    srv.heartBeatChan = make(chan time.Time, 100)
    srv.raft.toExecChan = srv.toExecChan
    srv.raft.heartBeatChan = srv.heartBeatChan
    go srv.execHandler()
    go srv.heartBeatHandler()
}

func (srv *Server) followerShutDown() {
    srv.raft.toExecChan = nil
    srv.raft.heartBeatChan = nil
    close(srv.toExecChan)
    close(srv.heartBeatChan)
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

func (srv *Server) UserLogout(token string) (bool, string) {
    ok := srv.validateAuth(token)
    if ok {
        user := srv.tokens[token]
        srv.deleteUserToken(user)
        return true, srv.messages.NoError
    }
    return false, srv.messages.UnrecognizedToken
}

func (srv *Server) UserLogin(username string, password string) (bool, string, string) {
    ok, err := srv.validateUser(username, password)
    if ok {
        user := srv.users[username]
        srv.generateUserToken(user)
        return true, srv.messages.NoError, user.token
    }
    return false, err.Error(), srv.messages.EmptyToken
}


func (srv *Server) RegisterUser(username string, password string) (bool, string, string) {
    if ok, err := srv.validateUserNameAndPassFormat(username, password); !ok {
        return ok, err.Error(), srv.messages.EmptyToken
    }
    if _, ok := srv.users[username]; ok {
        return false, srv.messages.UserAlreadyExist, srv.messages.EmptyToken
    }
    newUser := &User{
        Username:username,
        Password:password,
    }
    newUser.Init()
    srv.users[username] = newUser
    srv.generateUserToken(newUser)
    return true, srv.messages.NoError, newUser.token
}

func (srv *Server) DeleteUser(token string) (bool, string) {
    ok := srv.validateAuth(token)
    if ok {
        user := srv.tokens[token]
        srv.removeUser(user)
        return true, srv.messages.NoError
    }
    return false, srv.messages.UnrecognizedToken
}

func (srv *Server) Post(token string, content string) (bool, string) {
    ok := srv.validateAuth(token)
    if ok {
        user := srv.tokens[token]
        user.Post(content)
        return true, srv.messages.NoError
    }
    return false, srv.messages.UnrecognizedToken
}

func (srv *Server) Follow(token string, username string) (bool, string) {
    ok := srv.validateAuth(token)
    if ok {
        follower := srv.tokens[token]
        if user, hasUser := srv.users[username]; hasUser {
            follower.Follow(user)
            return true, srv.messages.NoError
        }
        return false, username + " " + srv.messages.UserNotExist
    }
    return false, srv.messages.UnrecognizedToken
}

func (srv *Server) UnFollow(token string, username string) (bool, string) {
    ok := srv.validateAuth(token)
    if ok {
        follower := srv.tokens[token]
        if user, hasUser := srv.users[username]; hasUser {
            follower.UnFollow(user)
            return true, srv.messages.NoError
        }
        return false, username + " " + srv.messages.UserNotExist
    }
    return false, srv.messages.UnrecognizedToken
}

func (srv *Server) GetMyContent(token string) (bool, string, []*cmd.Article) {
    ok := srv.validateAuth(token)
    if ok {
        user := srv.tokens[token]
        return true, srv.messages.NoError, user.GetMyContent()
    }
    return false, srv.messages.UnrecognizedToken, nil
}

func (srv *Server) GetFollower(token string) (bool, string, []*cmd.Relationship) {
    ok := srv.validateAuth(token)
    if ok {
        relationships := make([]*cmd.Relationship, 0, len(srv.users))
        follower := srv.tokens[token]
        for username, _ := range srv.users {
            if username != follower.Username {
                _, isFollowing := follower.following[username]
                relationships = append(relationships, &cmd.Relationship{username, isFollowing})
            }
        }
        return true, srv.messages.NoError, relationships
    }
    return false, srv.messages.UnrecognizedToken, nil
}

func (srv *Server) getFuncAndParameters(cmdValue reflect.Value) (reflect.Value, []reflect.Value) {
    srvValue := reflect.ValueOf(srv)
    f := srvValue.MethodByName(cmdValue.Type().Name())
    if !f.IsValid() {
            return f, nil
    }

    numberOfParameters := cmdValue.Field(0).Elem().NumField()
    parameters := make([]reflect.Value, numberOfParameters, numberOfParameters)
    for i:=0; i<numberOfParameters; i++ {
        parameters[i] = cmdValue.Field(0).Elem().Field(i)
    }

    return f, parameters
}

func (srv *Server) isReadOnly(cmdValue reflect.Value) bool {
    switch (cmdValue.Type().Name()) {
        case "GetMyFollower", "GetMyContent":
            return true
        default:
            return false
    }
}

func (srv *Server) exec(cmdValue reflect.Value) {
    f, parameters := srv.getFuncAndParameters(cmdValue)
    replyType := cmdValue.Field(1).Type().Elem().Elem()
    reply := reflect.New(replyType)

    if f.IsValid() {
        if srv.isReadOnly(cmdValue) {
            srv.rwLock.RLock()
            defer srv.rwLock.RUnlock()
        } else {
            srv.rwLock.Lock()
            defer srv.rwLock.Unlock();
        }
        results := f.Call(parameters)
        for index, value := range results {
            reply.Elem().Field(index).Set(value)
        }
    } else {
        reply.Elem().Field(0).Set(reflect.ValueOf(false))
        reply.Elem().Field(1).Set(reflect.ValueOf(srv.messages.FunctionNotImplement))
    }

    fmt.Printf("Command: %v, reply: %v\n", cmdValue.Type().Name(), reply)
    cmdValue.Field(1).Send(reply)
}

func (srv *Server) replyNotLeader(cmdValue reflect.Value) {
    replyType := cmdValue.Field(1).Type().Elem().Elem()
    reply := reflect.New(replyType)
    reply.Elem().Field(0).Set(reflect.ValueOf(false))
    leaderAddress := ""
    if srv.raft.leaderId >= 0 {
        leaderAddress = srv.addressBook[srv.raft.leaderId]
    }
    reply.Elem().Field(1).Set(reflect.ValueOf(srv.messages.NotLeader + leaderAddress))
    cmdValue.Field(1).Send(reply)
}

func (srv *Server) commitHandler() {
    indexCount := make(map[int]int)
    for index := range(srv.commitChan) {
        if srv.raft.commitIndex > index {
            continue
        }

        if _, ok := indexCount[index]; !ok {
            indexCount[index] = 1
        }
        indexCount[index] = indexCount[index] + 1

        if indexCount[index] > srv.getMajority() {
            for commitIndex := srv.raft.commitIndex + 1; commitIndex <= index; commitIndex++ {
                fmt.Printf("Exec %v\n", srv.commandLogs)
                delete(indexCount, commitIndex)
                srv.exec(srv.commandLogs[commitIndex])
                delete(srv.commandLogs, commitIndex)
                srv.raft.commitIndex = commitIndex
            }
        }
    }
}

func (srv *Server) updateLastBeat() {
    for beatTime := range srv.raft.heartBeatChan {
        fmt.Printf("logTerms: %v\n", srv.raft.logTerms)
        srv.lastBeatTime = beatTime
    }
}

func (srv *Server) startVote() bool {
    count := 1
    srv.raft.term = srv.raft.term + 1
    srv.raft.voteFor = srv.id
    countChan := make(chan int, len(srv.addressBook))
    for index, _ := range srv.addressBook {
        if srv.id == index {
            continue
        }
        go func(index int) {
            client := RaftClient{}
            err := client.InitOnce(srv.network, srv.addressBook[index])
            if err != nil {
                countChan <- 0
                return
            }
            lastLogIndex, lastLogTerm := srv.raft.getLastIndexAndTerm()
            reply, err := client.RequestVote(
                srv.raft.term,
                srv.id,
                lastLogIndex,
                lastLogTerm)
            if err == nil && reply.VoteGranted {
                countChan <- 1
            } else {
                countChan <- 0
            }
        }(index)
    }
    times := 1
    for result := range countChan {
        times++
        count += result
        if count > srv.getMajority() || times == len(srv.addressBook) {
            break
        }
    }
    if count > srv.getMajority() {
        return true
    }
    return false
}

func (srv *Server) heartBeatHandler(){
    go srv.updateLastBeat()
    for {
        time.Sleep(srv.timeout)
        randomTimeout := time.Duration(mrand.Intn(3) + 2) * srv.timeout
        if time.Now().Sub(srv.lastBeatTime) > randomTimeout {
            srv.lastBeatTime = time.Now()
            fmt.Print("Leader timeout\n")
            electionTimer := 10 * srv.timeout
            startVoteChan := make(chan bool, 1)
            go func(){
                startVoteChan <- srv.startVote()
            }()
            select {
                case voteRes := <-startVoteChan:
                    fmt.Printf("Election result: %v\n", voteRes)
                    if voteRes {
                        srv.followerShutDown()
                        srv.leaderInit()
                        return
                    }
                case <-time.After(electionTimer):
                    fmt.Println("election timeout")
            }
        }
    }
}

func (srv *Server) execHandler() {
    srvValue := reflect.ValueOf(srv)
    for execID :=  range srv.raft.toExecChan {
        encodedCmd := srv.raft.logs[execID]
        funcName, parameters := srv.cmdFactory.Decode(encodedCmd)
        f := srvValue.MethodByName(funcName)
        f.Call(parameters)
        fmt.Print("Replicate execute command " + funcName + "\n")
    }
}

func (srv *Server) followerHandler(index int) {
    fmt.Print("Start to Connect with " + srv.addressBook[index] + "\n")
    client := RaftClient{address:srv.addressBook[index]}
    client.Init(srv.network, srv.addressBook[index])
    fmt.Print("Successfully Connect with " + srv.addressBook[index] + "\n")
    //delay := 1
    for {
        if !srv.raft.isLeader {
            break
        }
        nextIndex := srv.nextIndexs[index]
        var command string

        if srv.raft.index < nextIndex {
            command = ""
            nextIndex = srv.raft.index + 1
            time.Sleep(srv.timeout)
            //delay++
        } else {
            command = srv.raft.logs[nextIndex]
        }

        _, commandTerm := srv.raft.getIndexAndTerm(nextIndex)
        preLogIndex, preLogTerm := srv.raft.getIndexAndTerm(nextIndex - 1)

        reply, err := client.AppendEntry(
            srv.raft.term,
            srv.id,
            preLogIndex,
            preLogTerm,
            command,
            commandTerm,
            srv.raft.commitIndex,
        )
        fmt.Printf("Replicate to index:%v, term:%v, prevLogIndex:%v, preLogTerm:%v, command:%v, commit:%v\n", index, srv.raft.term, preLogIndex, preLogTerm, command, srv.raft.commitIndex)

        if err != nil {
            fmt.Print(err)
            client.Init(srv.network, srv.addressBook[index])
            continue
        }
        if reply.Term > srv.raft.term {
            srv.raft.toFollowerChan <- reply.Term
            break
        }
        if reply.Success {
            if command != "" {
                srv.nextIndexs[index]++
                if commandTerm == srv.raft.term {
                    srv.commitChan <- nextIndex
                }
            }
        } else {
            fmt.Printf("Replicate fail with node %v and index %v, nextIndexs %v\n", index, nextIndex, srv.nextIndexs)
            srv.nextIndexs[index]--
        }
    }
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

func (srv *Server) runCommands() {
    for cmd := range srv.commands {
        if !srv.raft.isLeader {
            srv.replyNotLeader(cmd)
            continue
        }
        if srv.isReadOnly(cmd) {
            go srv.exec(cmd)
        } else {
            encodedCmd := srv.cmdFactory.Encode(cmd)
            srv.raft.appendCommand(encodedCmd, srv.raft.term)
            srv.commandLogs[srv.raft.index] = cmd
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

