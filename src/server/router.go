package server

import (
    "net/http"
    "log"
    "regexp"
)

var validPath = regexp.MustCompile("^/(edit|save|view|login|loginresult|signup)/")

func makeHandler(fn func(http.ResponseWriter, *http.Request, *Server), srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, srv)
	}
}

func Route(srv *Server) {
	http.HandleFunc("/login/", makeHandler(LoginHandler, srv))
	http.HandleFunc("/loginresult/", makeHandler(LoginresultHandler, srv))
	http.HandleFunc("/signup/", makeHandler(SignupHandler, srv))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
