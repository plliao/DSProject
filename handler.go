package handler

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type LogResult struct {
	Name string
	Password string
	Result string
	Message string
}

func loginHandler(w http.ResponseWriter, r *http.Request, title string) {
	//http.Redirect(w, r, "/logresult", http.StatusFound)
	var u User
	renderTemplate(w, "login", u)
	u.Name = r.FormValue("name")
	u.Password = r.FormValue("password")
}
func loginresultHandler(w http.ResponseWriter, r *http.Request, title string) {
	logres := LogResult { Name:r.FormValue("name"), Password:r.FormValue("password")}
	if(r.FormValue("choose") == "Log in"){
		pw, ok := users[r.FormValue("name")]
		if(pw == r.FormValue("password")){
			logres.Result = "successfully"
			//
		}else{
			logres.Result = "failed"
			logres.Message = "Wrong user."
			if ok {
				logres.Message = "Wrong password."
			}
		}
		renderTemplate(w, "loginresult", logres)

	}else{
		tmp := User{Name: r.FormValue("name"), Password:r.FormValue("password")}
		users[r.FormValue("name")] = &tmp
		renderTemplate(w, "signup", logres)
		//http.Redirect(w, r, "/signup/", http.StatusFound)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request, title string) {
	new_user := User{ Name:r.FormValue("name"), Password:r.FormValue("password")} 
	renderTemplate(w, "signup", new_user)
}
