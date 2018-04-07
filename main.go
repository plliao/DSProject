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

func main() {
    port := flag.Int("port", 8080, "Serving port")
    pagesDir := flag.String("d", "pages", "Default directory of HTML pages")

    htmls := []string{
        "login",
        "loginresult",
        "signup",
    }

    var srv server.Server
    srv.Init()
    registerHTMLs(&srv, htmls, *pagesDir)
    srv.Start(strconv.Itoa(*port))
}
