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

