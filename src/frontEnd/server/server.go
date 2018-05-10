package server

import (
    "html/template"
    "net/http"
    "log"
    "net/rpc"
    "fmt"
    "bufio"
    "os"
)

type Server struct {
    htmls map[string]string // name -> filepath
    handlers map[string]http.HandlerFunc // api -> handler
    templates *template.Template
    serverAddress []string
    leaderId int
    leaderAddress string
    network string
}

func (srv *Server) GetConnectInfo() (string, string) {
    return srv.leaderAddress, srv.network
}

func (srv *Server) GetAddressBook() {
    return srv.serverAddress
}

func (srv *Server) SetConnectInfo(address string, network string) {
    srv.leaderAddress = address
    srv.network = network
}

func (srv *Server) TryNextAddress() {
    srv.leaderId = (srv.leaderId + 1) % len(srv.serverAddress)
    srv.leaderAddress = srv.serverAddress[srv.leaderId]
}

func (srv *Server)ClientConnect() (*rpc.Client, error){
    var err error
    var client *rpc.Client
    for i := 0; i< len(srv.serverAddress); i++{
        client, err := rpc.DialHTTP(srv.network, srv.serverAddress[i])
        if(err == nil){
            return client, err
        }
    }
    return client, err
}

func (srv *Server) InitialDial(network string, filePath string){
    file, err := os.Open(filePath)
    srv.network = network
    if err != nil {
        return
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        srv.serverAddress = append(srv.serverAddress, scanner.Text())
    }
    srv.leaderId = 0
    srv.leaderAddress = srv.serverAddress[srv.leaderId]
}

func (srv *Server) Init() {
    srv.htmls = make(map[string]string)
    srv.handlers = make(map[string]http.HandlerFunc)
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

