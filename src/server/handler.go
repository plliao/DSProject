package server

import (
	"net/http"
    "net/url"
    "errors"
    "fmt"
)

type ServerHandlerFunc func(http.ResponseWriter, *http.Request, *Server)

type HandlerFuncFactory struct {
}

func (factory *HandlerFuncFactory) CreateByServerHandlerFunc(
        serverHandler ServerHandlerFunc, srv *Server) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        serverHandler(w, r, srv)
    }
}

type LoginPage struct {
	Message string
}

type FollowButton struct {
    Name string
    Action string
    User *User
}

type FollowPage struct {
    User *User
    Auth string
    FollowList []FollowButton
}

func LoginHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
    message := r.URL.Query().Get("message")
	  RenderTemplate(w, srv, "login", LoginPage{message})
}

func HomeHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
    token := r.FormValue("Auth")
    var user *User
    if srv.ValidateAuth(token) {
        user = srv.tokens[token]
        post := r.FormValue("article")
        if post != "" {
            user.Post(post)
        }
        logout := r.FormValue("logout")
        if logout != "" {
            srv.UserLogout(user)
            http.Redirect(w, r, "/login/", http.StatusFound)
            return
        }
    } else {
        username := r.FormValue("name")
        password := r.FormValue("password")
        choose := r.FormValue("choose")

        ok := false
        var err error
        if(choose == "Log in"){
            ok, err = srv.ValidateUser(username, password)
        } else if (choose == "Sign up") {
            ok, err = srv.RegisterUser(username, password)
        } else {
            err = errors.New("")
        }

        if !ok {
            loginURLValues := url.Values{}
            loginURLValues.Set("message", err.Error())
            http.Redirect(w, r, "/login/?" + loginURLValues.Encode(), http.StatusFound)
            return
        }
        user = srv.users[username]
        srv.UserLogin(user)
    }
	RenderTemplate(w, srv, "home", user)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request, srv *Server){
    token := r.FormValue("Auth")
    fmt.Printf("token: %s\n", token)
    var user *User
    var follow_action []FollowButton
    var follow_ FollowPage
    if srv.ValidateAuth(token) {
        user = srv.tokens[token]

        if(r.FormValue("FollowOrNot") == "Unfollow"){
            user.UnFollow(srv.users[r.FormValue("Target")])
        }else if(r.FormValue("FollowOrNot") == "Follow"){
            user.Follow(srv.users[r.FormValue("Target")])
        }

        for u := range(srv.users){
            _, ok := user.following[u]
            if(u != user.Username){
                if(ok){
                    follow_action = append(follow_action, FollowButton{Name:u, Action:"Unfollow", User:user})
                }else{
                        follow_action = append(follow_action, FollowButton{Name:u, Action:"Follow", User:user})
                }
            }
        }
        follow_ = FollowPage{ User: user, FollowList: follow_action}
    } else {
        /*username := r.FormValue("name")
        password := r.FormValue("password")
        choose := r.FormValue("choose")

        ok := false
        var err error
        if(choose == "Log in"){
            ok, err = srv.ValidateUser(username, password)
        } else if (choose == "Sign up") {
            ok, err = srv.RegisterUser(username, password)
        } else {
            err = errors.New("")
        }

        if !ok {
            loginURLValues := url.Values{}
            loginURLValues.Set("message", err.Error())
            http.Redirect(w, r, "/login/?" + loginURLValues.Encode(), http.StatusFound)
            return
        }
        user = srv.users[username]
        srv.UserLogin(user)*/
        http.Redirect(w, r, "/login/", http.StatusFound)
    }
    RenderTemplate(w, srv, "profile", follow_)
}

