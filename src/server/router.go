package server

import (
    "net/http"
    "regexp"
)

var validPath = regexp.MustCompile("^/(edit|save|view|login|loginresult|signup)/")

func validateURL(
        handlerFunc func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		handlerFunc(w, r)
	}
}

func Route(srv *Server) {
    var factory HandlerFuncFactory
	http.HandleFunc("/login/", validateURL(factory.CreateByServerHandlerFunc(LoginHandler, srv)))
	http.HandleFunc("/loginresult/", validateURL(factory.CreateByServerHandlerFunc(LoginresultHandler, srv)))
	http.HandleFunc("/signup/", validateURL(factory.CreateByServerHandlerFunc(SignupHandler, srv)))
}
