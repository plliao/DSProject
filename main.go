package main

import (
	//render "./libs/render"
	//handler "./libs/handler"
    "server"
	//"router"
	//"server"
)

func main() {
	var srv server.Server
	srv.Users = make(map[string]*server.User)
	server.Route(&srv)
}
