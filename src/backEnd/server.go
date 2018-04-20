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
    tokens map[string]*User
}

func (srv *Server) Init() {
    srv.users = make(map[string]*User)
    srv.tokens = make(map[string]*User)
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

func (srv *Server) DeleteUser(user *User) {
    srv.UserLogout(user)
    for _, follower := range user.followers {
        follower.UnFollow(user)
    }
    delete(srv.users, user.Username)
}

func (srv *Server) Start(port string) {

}
