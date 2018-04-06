package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"server"
)

func main() {
	var server Server
	server.users = make(map[string]*User)
	route()
}
