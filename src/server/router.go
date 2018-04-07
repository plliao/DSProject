package router

import (
    "net/http"
    "log"
    "regexp"
    server "../server"
    handler "../handler"
)

var validPath = regexp.MustCompile("^/(edit|save|view|login|loginresult|signup)/")

func makeHandler(fn func(http.ResponseWriter, *http.Request, *server.Server), srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, srv)
	}
}

func Route(srv *server.Server) {
	http.HandleFunc("/login/", makeHandler(handler.LoginHandler, srv))
	http.HandleFunc("/loginresult/", makeHandler(handler.LoginresultHandler, srv))
	http.HandleFunc("/signup/", makeHandler(handler.SignupHandler, srv))
	log.Fatal(http.ListenAndServe(":8080", nil))
}