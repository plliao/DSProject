package main

import (
	//render "./libs/render"
	//handler "./libs/handler"
	user "./libs/user"
	server "./libs/server"
	router "./libs/router"
	//"router"
	//"server"
)

func main() {
	var server server.Server
	server.Users = make(map[string]*user.User)
	router.Route(&server)
}
