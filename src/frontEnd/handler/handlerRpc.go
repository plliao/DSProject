package handler

import (
	//"net/http"
    //"net/url"
    "frontEnd/server"
    "backEnd/cmd"
    //"errors"
    //"strings"
    //"html/template"
    //"log"
)

func ClientPostRPC(token string, post string, srv *server.Server) (error, cmd.PostReply){
	args := cmd.PostArgs { Token:token, Content:post }
	var reply cmd.PostReply
    err := srv.SrvClient.Call("Service." + "Post", args, &reply)
    return err, reply
}

func ClientLogoutRPC(token string, srv *server.Server) (error, cmd.UserLogoutReply){
	args := cmd.UserLogoutArgs{ Token:token }
    var reply cmd.UserLogoutReply
	err := srv.SrvClient.Call("Service." + "UserLogout", args, &reply)
	return err, reply
}