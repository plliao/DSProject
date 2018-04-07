package server

import (
	"net/http"
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

func LoginHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
	RenderTemplate(w, srv, "login", LoginPage{""})
}

func LoginresultHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
    username := r.FormValue("name")
    password := r.FormValue("password")

    var ok bool
    var err error
	if(r.FormValue("choose") == "Log in"){
        ok, err = srv.ValidateUser(username, password)
	} else {
        ok, err = srv.RegisterUser(username, password)
	}

    if ok {
        user := srv.users[username]
        user.Post("Articl1")
        user.Post("Articl2")
        user.Post("Articl3")
	    RenderTemplate(w, srv, "home", user)
    } else {
	    RenderTemplate(w, srv, "login", LoginPage{err.Error()})
    }
}

func HomeHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
    username := r.FormValue("name")
    user, ok := srv.users[username]
    fmt.Printf("Home%+v\n", user)
    if !ok {
        return
    } else {
	    RenderTemplate(w, srv, "home", user)
    }
}
