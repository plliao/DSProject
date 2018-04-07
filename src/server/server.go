package server

import (
    "html/template"
    "net/http"
    "log"
)

type Server struct {
    users map[string]*User
    htmls map[string]string // name -> filepath
    templates *template.Template

}

func (srv *Server) Init() {
    srv.users = make(map[string]*User)
    srv.htmls = make(map[string]string)

}

func (srv *Server) RegisterHTML(name string, path string) {
    srv.htmls[name] = path
}

func (srv *Server) Start(port string) {
    srv.createTemplates()
    Route(srv)
    log.Fatal(http.ListenAndServe(":" + port, nil))
}

func (srv *Server) createTemplates() {
    filepaths := make([]string, 0, len(srv.htmls))
    for _, filepath := range srv.htmls {
        filepaths = append(filepaths, filepath)
    }
    srv.templates = CreateTemplates(filepaths...)
}
