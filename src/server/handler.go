package server

import (
	"net/http"
    "net/url"
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
    username := r.FormValue("name")
    password := r.FormValue("password")
    choose := r.FormValue("choose")

    ok := false
    var err error
	if(choose == "Log in"){
        ok, err = srv.ValidateUser(username, password)
	} else if (choose == "Sign up") {
        ok, err = srv.RegisterUser(username, password)
	}

    if ok {
        user := srv.users[username]
        user.Post("Articl1")
	    RenderTemplate(w, srv, "home", user)
    } else {
        loginURLValues := url.Values{}
        loginURLValues.Set("message", err.Error())
        http.Redirect(w, r, "/login/?" + loginURLValues.Encode(), http.StatusFound)
    }
}

