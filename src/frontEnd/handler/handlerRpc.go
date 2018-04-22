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

func ClientRegisterUserRPC(username string, password string, srv *server.Server) (error, ClientReply){
	args := cmd.RegisterUserArgs{ Username:username, Password:password }
    var reply cmd.RegisterUserReply
    err := srv.SrvClient.Call("Service." + "RegisterUser", args, &reply)
    return err, ClientReply{ Ok:reply.Ok, Token:reply.Token, Error:reply.Error}
}

func ClientUserLoginRPC(username string, password string, srv *server.Server) (error, ClientReply){
	args := cmd.UserLoginArgs{ Username:username, Password:password }
    var reply cmd.UserLoginReply
    err := srv.SrvClient.Call("Service." + "UserLogin", args, &reply)
    return err, ClientReply{ Ok:reply.Ok, Token:reply.Token, Error:reply.Error}
}

func ClientGetMyContentRPC(token string, srv *server.Server) (error, cmd.GetMyContentReply){
	args := cmd.GetMyContentArgs{ Token:token }
    var reply cmd.GetMyContentReply
    err := srv.SrvClient.Call("Service." + "GetMyContent", args, &reply)
    return err, reply
}

func ClientGetFollowerRPC(token string, srv *server.Server) (string, cmd.GetFollowerReply){
	args := cmd.GetFollowerArgs{ Token:token }
    var reply cmd.GetFollowerReply 
    err := srv.SrvClient.Call("Service." + "GetFollower", args, &reply)
    errmsg := ""
    if(err != nil){
    	errmsg = "Authentication failed."
    }
    return errmsg, reply
}

func ClientDeleteUserRPC(token string, srv *server.Server) (string, cmd.DeleteUserReply){
	args := cmd.DeleteUserArgs{ Token:token }
    var reply cmd.DeleteUserReply
    err := srv.SrvClient.Call("Service." + "DeleteUser", args, &reply)
    errmsg := ""
    if(err != nil || !reply.Ok){
        errmsg = "Delete Failed. Please log in again."
    }
    return errmsg, reply
}

func ClientUnFollowRPC(token string, target string, srv *server.Server) (error, cmd.UnFollowReply){
	args := cmd.UnFollowArgs{ Token:token, Username:target}
    var reply cmd.UnFollowReply
    err := srv.SrvClient.Call("Service." + "UnFollow", args, &reply)
    return err, reply
}

func ClientFollowRPC(token string, target string, srv *server.Server) (error, cmd.FollowReply){
	args := cmd.FollowArgs{ Token:token, Username:target}
    var reply cmd.FollowReply
    err := srv.SrvClient.Call("Service." + "Follow", args, &reply)
    return err, reply
}