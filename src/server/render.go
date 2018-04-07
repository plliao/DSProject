package server

import (
    "html/template"
    "net/http"
    "path"
)

func CreateTemplates(filepaths ... string) *template.Template {
    return template.Must(template.ParseFiles(filepaths...))
}

func RenderTemplate(w http.ResponseWriter, srv *Server, name string, data interface{}) {
    template_name := path.Base(srv.htmls[name])
	err := srv.templates.ExecuteTemplate(w, template_name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
