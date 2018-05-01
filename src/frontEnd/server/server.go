package server

import (
    "html/template"
    "net/http"
    "log"
    "net/rpc"
    "fmt"
)

type Server struct {
    htmls map[string]string // name -> filepath
    handlers map[string]http.HandlerFunc // api -> handler
    templates *template.Template
    serverAddress string
    network string
}

func (srv *Server)ClientConnect() (*rpc.Client, error){
    client, err := rpc.DialHTTP(srv.network, srv.serverAddress)
    return client, err
}

func (srv *Server) InitialDial(network string, serverAddress string){
    srv.serverAddress = serverAddress
    srv.network = network
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
    fmt.Print("FrontEnd server Start...\n")
    log.Fatal(http.ListenAndServe(":" + port, nil))
}

func (srv *Server) createTemplates() {
    filepaths := make([]string, 0, len(srv.htmls))
    for _, filepath := range srv.htmls {
        filepaths = append(filepaths, filepath)
    }
    srv.templates = CreateTemplates(filepaths...)
}

func (srv *Server) GetAPIAndHandlers() ([]string, []http.HandlerFunc) {
    apis := make([]string, 0, len(srv.handlers))
    handlers := make([]http.HandlerFunc, 0, len(srv.handlers))
    for api, handler := range srv.handlers {
        apis = append(apis, api)
        handlers = append(handlers, handler)
    }
    return apis, handlers
}

