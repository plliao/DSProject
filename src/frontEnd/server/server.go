package server

import (
    "html/template"
    "net/http"
    "log"
    "net/rpc"
)

type Server struct {
    htmls map[string]string // name -> filepath
    handlers map[string]http.HandlerFunc // api -> handler
    templates *template.Template

    SrvClient *rpc.Client
}

func (srv *Server) Init() {
    srv.htmls = make(map[string]string)
    srv.handlers = make(map[string]http.HandlerFunc)
    //srv.client = make(rpc.Client)
}

func (srv *Server) RegisterHTML(name string, path string) {
    srv.htmls[name] = path
}

func (srv *Server) RegisterHandlerFunc(api string, handler http.HandlerFunc) {
    srv.handlers[api] = handler
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

func (srv *Server) GetAPI() []string {
    apis := make([]string, 0, len(srv.handlers))
    for api, _ := range srv.handlers {
        apis = append(apis, api)
    }
    return apis
}

func (srv *Server) GetHandlers() []http.HandlerFunc {
    handlers := make([]http.HandlerFunc, 0, len(srv.handlers))
    for _, handler := range srv.handlers {
        handlers = append(handlers, handler)
    }
    return handlers
}
