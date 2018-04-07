package main

import (
    "flag"
    "strconv"
    "server"
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

func main() {
    port := flag.Int("port", 8080, "Serving port")
    pagesDir := flag.String("d", "pages", "Default directory of HTML pages")
    flag.Parse()

    htmls := []string{
        "login",
        "loginresult",
        "signup",
    }

    apiToServerHandlerFuncMap := map[string]server.ServerHandlerFunc{
        "login":server.LoginHandler,
        "loginresult":server.LoginresultHandler,
        "signup":server.SignupHandler,
    }

    var srv server.Server
    srv.Init()
    registerHTMLs(&srv, htmls, *pagesDir)
    createAndRegisterServerHandlers(&srv, apiToServerHandlerFuncMap)
    srv.Start(strconv.Itoa(*port))
}
