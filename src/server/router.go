package server

import (
    "net/http"
    "regexp"
    "strings"
)

type Middleware interface {
    Decorate(http.HandlerFunc) http.HandlerFunc
}

type MiddlewareManager struct {
    decorators []*Middleware
}

func (manager *MiddlewareManager) Decorate(handlerFunc http.HandlerFunc) http.HandlerFunc {
    decoratedHandlerFunc := handlerFunc
    for _, decorator := range manager.decorators {
        decoratedHandlerFunc = (*decorator).Decorate(decoratedHandlerFunc)
    }
    return decoratedHandlerFunc
}

func (manager *MiddlewareManager) RegisterMiddleware(middleware Middleware) {
    manager.decorators = append(manager.decorators, &middleware)
}


type URLValidator struct {
    validPath *regexp.Regexp
}

func (validator URLValidator) Decorate(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validator.validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		handlerFunc(w, r)
	}
}

func Route(srv *Server) {
    var manager MiddlewareManager

    apis := srv.GetAPI()
    validator := URLValidator{
        validPath:regexp.MustCompile("^/" + strings.Join(apis, "|") + "/"),
    }
    manager.RegisterMiddleware(validator)

    handlers := srv.GetHandlers()
    for index, handler := range handlers {
        decoratedHandler := manager.Decorate(handler)
        http.HandleFunc("/" + apis[index] + "/", decoratedHandler)
    }
}
