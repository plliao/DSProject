package render

import {
    "html/template"
    "net/http"
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "login.html", "loginresult.html", "signup.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}