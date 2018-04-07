package server

import (
    "html/template"
    "net/http"
)

var templates = template.Must(template.ParseFiles("pages/login.html", "pages/loginresult.html", "pages/signup.html"))

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
