package handler

import (
	"net/http"
	server "../server"
	user "../user"
	render "../render"
)

type LogResult struct {
	Name string
	Password string
	Result string
	Message string
}

func LoginHandler(w http.ResponseWriter, r *http.Request, srv *server.Server) {
	//http.Redirect(w, r, "/logresult", http.StatusFound)
	var u user.User
	render.RenderTemplate(w, "login", u)
	u.Username = r.FormValue("name")
	u.Password = r.FormValue("password")
}
func LoginresultHandler(w http.ResponseWriter, r *http.Request, srv *server.Server) {
	logres := LogResult { Name:r.FormValue("name"), Password:r.FormValue("password")}
	if(r.FormValue("choose") == "Log in"){
		pw, ok := srv.Users[r.FormValue("name")]
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
		render.RenderTemplate(w, "loginresult", logres)

	}else{
		tmp := user.User{Username: r.FormValue("name"), Password:r.FormValue("password")}
		srv.Users[r.FormValue("name")] = &tmp
		render.RenderTemplate(w, "signup", logres)
		//http.Redirect(w, r, "/signup/", http.StatusFound)
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request, srv *server.Server) {
	new_user := user.User{ Username:r.FormValue("name"), Password:r.FormValue("password")} 
	render.RenderTemplate(w, "signup", new_user)
}
