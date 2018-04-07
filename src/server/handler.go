package server

import (
	"net/http"
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

type LogResult struct {
	Name string
	Password string
	Result string
	Message string
}


func LoginHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
	//http.Redirect(w, r, "/logresult", http.StatusFound)
	var u User
	RenderTemplate(w, srv, "login", u)
	u.Username = r.FormValue("name")
	u.Password = r.FormValue("password")
}

func LoginresultHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
	logres := LogResult { Name:r.FormValue("name"), Password:r.FormValue("password")}
	if(r.FormValue("choose") == "Log in"){
		pw, ok := srv.users[r.FormValue("name")]
		if(pw.Password == r.FormValue("password")){
			logres.Result = "successfully"
			//
		}else{
			logres.Result = "failed"
			logres.Message = "Wrong user."
			if ok {
				logres.Message = "Wrong password."
			}
		}
		RenderTemplate(w, srv, "loginresult", logres)

	}else{
		tmp := User{Username: r.FormValue("name"), Password:r.FormValue("password")}
		srv.users[r.FormValue("name")] = &tmp
		RenderTemplate(w, srv, "signup", logres)
		//http.Redirect(w, r, "/signup/", http.StatusFound)
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
	new_user := User{ Username:r.FormValue("name"), Password:r.FormValue("password")}
	RenderTemplate(w, srv, "signup", new_user)
}
