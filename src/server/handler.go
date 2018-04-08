package server

import (
	  "net/http"
    "net/url"
    "errors"
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

type FollowPage struct {
    Following []string
    Username string
}

func LoginHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
    message := r.URL.Query().Get("message")
	  RenderTemplate(w, srv, "login", LoginPage{message})
}

func HomeHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
    token := r.FormValue("auth")
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

func ProfileHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
    username := r.FormValue("name")
    u := FollowPage{Username:username}
    for _, fu := range(srv.users[username].following){
        u.Following = append(u.Following, fu.Username)
    }
    RenderTemplate(w, srv, "profile", u)
}
