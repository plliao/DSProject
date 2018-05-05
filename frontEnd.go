package main

import (
    "flag"
    "strconv"
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

//func createBackEndAddress(srv *server.Server, network string, serverAddress string) {
//    srv.InitialDial(network, serverAddress)
//    return
//}

func main() {
    httpPort := flag.Int("port", 8877, "Serving port")
    //backendServer := flag.String("b", "71.125.15.10:80", "backend server ip and port e.g. 127.0.0.1:80")
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
        "profile":handler.ProfileHandler,
    }

    var srv server.Server
    srv.InitialDial("tcp", "config.txt")
    srv.Init()
    registerHTMLs(&srv, htmls, *pagesDir)
    createAndRegisterServerHandlers(&srv, apiToServerHandlerFuncMap)
    //serverAddress := *backendServer
    //createBackEndAddress(&srv, "tcp", serverAddress)
    srv.Start(strconv.Itoa(*httpPort))
}
