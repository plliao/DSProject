package server

import (
    "html/template"
    "net/http"
    "log"
    "regexp"
    "errors"
    "crypto/rand"
    "fmt"
)

type Server struct {
    users map[string]*User
    htmls map[string]string // name -> filepath
    handlers map[string]http.HandlerFunc // api -> handler
    templates *template.Template

    validUserName *regexp.Regexp
    validPassword *regexp.Regexp

    tokens map[string]*User
}

func (srv *Server) Init() {
    srv.users = make(map[string]*User)
    srv.htmls = make(map[string]string)
    srv.tokens = make(map[string]*User)
    srv.handlers = make(map[string]http.HandlerFunc)

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

func (srv *Server) RegisterUser(username string, password string) (bool, error) {
    if ok, err := srv.validateUserNameAndPassFormat(username, password); !ok {
        return ok, err
    }
    if _, ok := srv.users[username]; ok {
        return false, errors.New("User already exists")
    }
    newUser := &User{
        Username:username,
        Password:password,
    }
    newUser.Init()
    srv.users[username] = newUser
    return true, nil
}

func (srv *Server) RegisterHTML(name string, path string) {
    srv.htmls[name] = path
}

func (srv *Server) RegisterHandlerFunc(api string, handler http.HandlerFunc) {
    srv.handlers[api] = handler
}

func (srv *Server) Start(port string) {
    srv.createTemplates()
    Route(srv)
    log.Fatal(http.ListenAndServe(":" + port, nil))
}

func (srv *Server) createTemplates() {
    filepaths := make([]string, 0, len(srv.htmls))
    for _, filepath := range srv.htmls {
        filepaths = append(filepaths, filepath)
    }
    srv.templates = CreateTemplates(filepaths...)
}

func (srv *Server) GetAPI() []string {
    apis := make([]string, 0, len(srv.handlers))
    for api, _ := range srv.handlers {
        apis = append(apis, api)
    }
    return apis
}

func (srv *Server) GetHandlers() []http.HandlerFunc {
    handlers := make([]http.HandlerFunc, 0, len(srv.handlers))
    for _, handler := range srv.handlers {
        handlers = append(handlers, handler)
    }
    return handlers
}
