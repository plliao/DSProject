package backEnd

import (
    "net"
    "net/http"
    "net/rpc"
    "log"
    "regexp"
    "errors"
    "crypto/rand"
    "fmt"
    "reflect"
    "backEnd/cmd"
    "sync"
)

type Server struct {
    users map[string]*User
    tokens map[string]*User
    commands chan reflect.Value

    rwLock *sync.RWMutex
    service *Service
    messages BackEndMessages

    validUserName *regexp.Regexp
    validPassword *regexp.Regexp

    addressBook []string
    raft Raft
    commitChan chan int
    nextIndexs []int
}

func (srv *Server) Init() {
    srv.users = make(map[string]*User)
    srv.tokens = make(map[string]*User)
    srv.commands = make(chan reflect.Value, 100)

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
    }

    srv.validUserName, _ = regexp.Compile("^[a-zA-Z0-9]{4,10}$")
    srv.validPassword, _ = regexp.Compile("^[a-zA-Z0-9]{4,10}$")

    raft = Raft{
        term:0,
        index:0,
        commitIndex:-1
    }
    addressBook = make([]string, 0)
}

func (srv *Server) RegisterAddress(address string, port string) {
    addressBook = append(addressBook, address + ":" + port)
    nextIndexs = append(nextIndexs, 0)
}

func (srv *Server) masterInit() {
    nextIndexs = make([]int, len(addressBook))
    commitChan = make(chan int, 100)
}

func (srv *Server) masterShutDown() {
    nextIndexs = nil
    close(commitChan)
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
    token := make([]byte, 6)
    rand.Read(token)
    user.token = fmt.Sprintf("%x", token)
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

func (srv *Server) commitHandler() {
    indexCount := make([int]int)
    for index := range(srv.commitChan) {
        if srv.raft.commitIndex > index {
            continue
        }

        if _, ok := indexCount[index]; !ok {
            indexCount[index] = 1
        }
        indexCount[index] = indexCount[index] + 1

        if indexCount[index] > srv.getMajority() {
            for commitIndex = srv.raft.commitIndex + 1; commitIndex <= index; commitIndex++ {
                delete(indexCount, commitIndex)
                srv.exec(srv.raft.logs[commitIndex])
                srv.raft.commitIndex = commitIndex
            }
        }
    }
}

func (srv *Server) slaveHandler(index int) {
    client := RaftClient{srv.addressBook[index]}
    client.Init()

}

func (srv *Server) runCommands() {
    for cmd := range srv.commands {
        go srv.exec(cmd)
    }
}

func (srv *Server) Start(port string) {
    go srv.runCommands()

    rpc.Register(srv.service)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":" + port)
    if e != nil {
        log.Fatal("listen error:", e)
    }
    http.Serve(l, nil)
}

