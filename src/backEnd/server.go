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
)

type Command interface {
}

type Server struct {
    users map[string]*User
    tokens map[string]*User
    commands chan reflect.Value

    service *Service

    validUserName *regexp.Regexp
    validPassword *regexp.Regexp
}

func (srv *Server) Init() {
    srv.users = make(map[string]*User)
    srv.tokens = make(map[string]*User)
    srv.commands = make(chan reflect.Value, 100)

    srv.service = &Service{srv.commands}

    srv.validUserName, _ = regexp.Compile("^[a-zA-Z0-9]{4,10}$")
    srv.validPassword, _ = regexp.Compile("^[a-zA-Z0-9]{4,10}$")
}

func (srv *Server) validateUserNameAndPassFormat(username string, password string) (bool, error) {
    if !srv.validUserName.Match([]byte(username)) || !srv.validPassword.Match([]byte(password)) {
        return false, errors.New("Incorrect username or password format")
    }
    return true, nil
}

func (srv *Server) ValidateUser(username string, password string) (bool, error) {
    if ok, err := srv.validateUserNameAndPassFormat(username, password); !ok {
        return ok, err
    }
    if user, ok := srv.users[username]; ok {
        if user.Password != password {
            return false, errors.New("Incorrect password")
        }
        return true, nil
    }
    return false, errors.New("User not exists")
}

func (srv *Server) ValidateAuth(token string) bool {
    if _, ok := srv.tokens[token]; ok {
        return true
    }
    return false
}

func (srv *Server) UserLogout(user *User) {
    if user.token != "" {
        delete(srv.tokens, user.token)
        user.token = ""
    }
}

func (srv *Server) UserLogin(user *User) {
    srv.UserLogout(user)
    token := make([]byte, 6)
    rand.Read(token)
    user.token = fmt.Sprintf("%x", token)
    srv.tokens[user.token] = user
}

func (srv *Server) RegisterUser(username string, password string) (bool, string, error) {
    if ok, err := srv.validateUserNameAndPassFormat(username, password); !ok {
        return ok, "", err
    }
    if _, ok := srv.users[username]; ok {
        return false, "", errors.New("User already exists")
    }
    newUser := &User{
        Username:username,
        Password:password,
    }
    newUser.Init()
    srv.users[username] = newUser
    srv.UserLogin(newUser)
    return true, newUser.token, nil
}

func (srv *Server) DeleteUser(user *User) {
    srv.UserLogout(user)
    for _, follower := range user.followers {
        follower.UnFollow(user)
    }
    delete(srv.users, user.Username)
}

func (srv *Server) getFuncAndParameters(cmdValue reflect.Value) (reflect.Value, []reflect.Value) {
    offset := 0

    srvValue := reflect.ValueOf(srv)
    f := srvValue.MethodByName(cmdValue.Type().Name())
    if !f.IsValid() {
        token := cmdValue.Field(0).Elem().Field(0).Interface().(string)
        user := srv.tokens[token]
        userValue := reflect.ValueOf(user)
        f = userValue.MethodByName(cmdValue.Type().Name())
        offset = 1
    }

    numberOfParameters := cmdValue.Field(0).Elem().NumField() - offset
    parameters := make([]reflect.Value, numberOfParameters, numberOfParameters)
    for i:=0; i<numberOfParameters; i++ {
        parameters[i] = cmdValue.Field(0).Elem().Field(i + offset)
    }

    return f, parameters
}

func (srv *Server) exec(cmdValue reflect.Value) {
    f, parameters := srv.getFuncAndParameters(cmdValue)
    fmt.Printf("+%v\n", f)
    results := f.Call(parameters)
    fmt.Printf("+%v\n", results)

    replyType := cmdValue.Field(1).Type().Elem().Elem()
    reply := reflect.New(replyType)
    offset := 0
    for index, value := range results {
        if index == 0 && value.Type().Name() != "bool" {
            reply.Elem().Field(0).Set(reflect.ValueOf(true))
            offset = 1
        }
        reply.Elem().Field(index + offset).Set(value)
    }
    fmt.Printf("+%v\n", reply)
    cmdValue.Field(1).Send(reply)
}

func (srv *Server) runCommands() {
    for cmd := range srv.commands {
        srv.exec(cmd)
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
