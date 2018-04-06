package router

import {
    "net/http"
    "handler"
    "log"
}

func route() {
	http.HandleFunc("/login/", makeHandler(loginHandler))
	http.HandleFunc("/loginresult/", makeHandler(loginresultHandler))
	http.HandleFunc("/signup/", makeHandler(signupHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var validPath = regexp.MustCompile("^/(edit|save|view|login|loginresult|signup)/")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[1])
	}
}

