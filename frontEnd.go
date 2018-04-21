package main

import (
    "flag"
    "strconv"
    "net/rpc"
    "frontEnd/server"
    "frontEnd/handler"
)

func registerHTMLs(srv *server.Server, htmls []string, pagesDir string) {
    surfix := ".html"
    for _, htmlName := range htmls {
        srv.RegisterHTML(htmlName, pagesDir + "/" + htmlName + surfix)
    }
}

func createAndRegisterServerHandlers(
        srv *server.Server, apiToServerHandlerFuncMap map[string]server.ServerHandlerFunc) {
    var factory server.HandlerFuncFactory
    for api, serverHandlerFunc := range apiToServerHandlerFuncMap {
        handlerFunc := factory.CreateByServerHandlerFunc(serverHandlerFunc, srv)
        srv.RegisterHandlerFunc(api, handlerFunc)
    }
}

func clientConnect(srv *server.Server, network string, serverAddress string) bool{
    var err error
    srv.SrvClient, err = rpc.DialHTTP(network, serverAddress)
    if(err == nil){
        return true
    }else{
        return false
    }
}

func main() {
    port := flag.Int("port", 8080, "Serving port")
    pagesDir := flag.String("d", "pages", "Default directory of HTML pages")
    flag.Parse()

    htmls := []string{
        "login",
        "home",
        "profile",
    }

    apiToServerHandlerFuncMap := map[string]server.ServerHandlerFunc{
        "login":handler.LoginHandler,
        "home":handler.HomeHandler,
        "profile": handler.ProfileHandler,
    }

    var srv server.Server
    srv.Init()
    registerHTMLs(&srv, htmls, *pagesDir)
    createAndRegisterServerHandlers(&srv, apiToServerHandlerFuncMap)
    serverAddress := "..."
    clientConnect(&srv, "tcp", serverAddress)
    srv.Start(strconv.Itoa(*port))
}
