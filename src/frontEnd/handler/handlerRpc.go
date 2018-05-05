package handler

import (
	//"net/http"
    //"net/url"
    "backEnd/cmd"
    "net/rpc"
    //"errors"
    //"strings"
    //"html/template"
    //"log"
)

func ClientCall(client *rpc.Client, service string, args interface, reply interface) (error, interface) {
    err := client.Call(service, args, &reply)
    if reply.Ok == false && reply.Error[:12] == "Not Leader: " {
        address := reply.Error[12:]
    }
}

func ClientPostRPC(token string, post string, client *rpc.Client) (error, cmd.PostReply){
	args := cmd.PostArgs { Token:token, Content:post }
	var reply cmd.PostReply
    err := client.Call("Service." + "Post", args, &reply)
    return err, reply
}

func ClientLogoutRPC(token string, client *rpc.Client) (error, cmd.UserLogoutReply){
	args := cmd.UserLogoutArgs{ Token:token }
    var reply cmd.UserLogoutReply
	err := client.Call("Service." + "UserLogout", args, &reply)
	return err, reply
}

func ClientRegisterUserRPC(username string, password string, client *rpc.Client) (error, ClientReply){
	args := cmd.RegisterUserArgs{ Username:username, Password:password }
    var reply cmd.RegisterUserReply
    err := client.Call("Service." + "RegisterUser", args, &reply)
    return err, ClientReply{ Ok:reply.Ok, Token:reply.Token, Error:reply.Error}
}

func ClientUserLoginRPC(username string, password string, client *rpc.Client) (error, ClientReply){
	args := cmd.UserLoginArgs{ Username:username, Password:password }
    var reply cmd.UserLoginReply
    err := client.Call("Service." + "UserLogin", args, &reply)
    return err, ClientReply{ Ok:reply.Ok, Token:reply.Token, Error:reply.Error}
}

func ClientGetMyContentRPC(token string, client *rpc.Client) (error, cmd.GetMyContentReply){
	args := cmd.GetMyContentArgs{ Token:token }
    var reply cmd.GetMyContentReply
    err := client.Call("Service." + "GetMyContent", args, &reply)
    return err, reply
}

func ClientGetFollowerRPC(token string, client *rpc.Client) (string, cmd.GetFollowerReply){
	args := cmd.GetFollowerArgs{ Token:token }
    var reply cmd.GetFollowerReply
    err := client.Call("Service." + "GetFollower", args, &reply)
    errmsg := ""
    if(err != nil){
    	errmsg = "Authentication failed."
    }
    return errmsg, reply
}

func ClientDeleteUserRPC(token string, client *rpc.Client) (string, cmd.DeleteUserReply){
	args := cmd.DeleteUserArgs{ Token:token }
    var reply cmd.DeleteUserReply
    err := client.Call("Service." + "DeleteUser", args, &reply)
    errmsg := ""
    if(err != nil || !reply.Ok){
        errmsg = "Delete Failed. Please log in again."
    }
    return errmsg, reply
}

func ClientUnFollowRPC(token string, target string, client *rpc.Client) (error, cmd.UnFollowReply){
	args := cmd.UnFollowArgs{ Token:token, Username:target}
    var reply cmd.UnFollowReply
    err := client.Call("Service." + "UnFollow", args, &reply)
    return err, reply
}

func ClientFollowRPC(token string, target string, client *rpc.Client) (error, cmd.FollowReply){
	args := cmd.FollowArgs{ Token:token, Username:target}
    var reply cmd.FollowReply
    err := client.Call("Service." + "Follow", args, &reply)
    return err, reply
}
