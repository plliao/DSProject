package handler

import (
	//"net/http"
    //"net/url"
    "backEnd/cmd"
    "net/rpc"
    "frontEnd/server"
    //"strings"
    //"html/template"
    //"log"
    //"fmt"
    "reflect"
)

func ClientCall(service string, args interface{}, replyType reflect.Type, srv *server.Server) (error, interface{}){
    i := 0
    addressBook, network := srv.GetConnectInfo()
    address := addressBook[i]
    for {
        client, errDial := rpc.DialHTTP(network, address)
        if(i >= len(addressBook)){
            i = 0
        }
        if(errDial == nil){
            reply := reflect.New(replyType)
            errRPC := client.Call(service, args, reply.Interface())
            ok := reply.Elem().Field(0).Interface().(bool)
            message := reply.Elem().Field(1).Interface().(string)
            NotLeader := ""
            //fmt.Print(len(message))
            if len(message) > 13{
                NotLeader = message[:11]
            //    fmt.Print("message[:11]", NotLeader)
            //    fmt.Print("\nmessage[:12]", message[:12])
            //    fmt.Print("\nmessage[12:]", message[12:])
            }
            if errRPC != nil {
                address = addressBook[i]
                i++
            } else if ok == false &&  NotLeader == "Not Leader: " {
                address = message[12:]
            } else {
                dupReply := reply.Interface()
                return errRPC, dupReply//reply.Interface()
            }
        } else{
            address = addressBook[i]
            i++
        }
    }
}

func ClientPostRPC(token string, post string, srv *server.Server) (error, cmd.PostReply){
	args := cmd.PostArgs { Token:token, Content:post }
    var reply cmd.PostReply
    err, replyInf := ClientCall("Service." + "Post", args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.PostReply)))
}

func ClientLogoutRPC(token string, srv *server.Server) (error, cmd.UserLogoutReply){
	args := cmd.UserLogoutArgs{ Token:token }
    var reply cmd.UserLogoutReply
	err, replyInf := ClientCall("Service." + "UserLogout", args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.UserLogoutReply)))
}

func ClientRegisterUserRPC(username string, password string, srv *server.Server) (error, ClientReply){
	args := cmd.RegisterUserArgs{ Username:username, Password:password }
    var reply ClientReply
    err, replyInf := ClientCall("Service." + "RegisterUser", args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*ClientReply))
}

func ClientUserLoginRPC(username string, password string, srv *server.Server) (error, ClientReply){
	args := cmd.UserLoginArgs{ Username:username, Password:password }
    var reply ClientReply
    err, replyInf := ClientCall("Service." + "UserLogin", args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*ClientReply))
}

func ClientGetMyContentRPC(token string, srv *server.Server) (error, cmd.GetMyContentReply){
	args := cmd.GetMyContentArgs{ Token:token }
    var reply cmd.GetMyContentReply
    err, replyInf := ClientCall("Service." + "GetMyContent", args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.GetMyContentReply)))
}

func ClientGetFollowerRPC(token string, srv *server.Server) (error, cmd.GetFollowerReply){
	args := cmd.GetFollowerArgs{ Token:token }
    var reply cmd.GetFollowerReply
    err, replyInf := ClientCall("Service." + "GetFollower", args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.GetFollowerReply)))
}

func ClientDeleteUserRPC(token string, srv *server.Server) (error, cmd.DeleteUserReply){
	args := cmd.DeleteUserArgs{ Token:token }
    var reply cmd.DeleteUserReply
    err, replyInf := ClientCall("Service." + "DeleteUser", args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.DeleteUserReply)))
}

func ClientUnFollowRPC(token string, target string, srv *server.Server) (error, cmd.UnFollowReply){
	args := cmd.UnFollowArgs{ Token:token, Username:target}
    var reply cmd.UnFollowReply
    err, replyInf := ClientCall("Service." + "UnFollow", args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.UnFollowReply)))
}

func ClientFollowRPC(token string, target string, srv *server.Server) (error, cmd.FollowReply){
	args := cmd.FollowArgs{ Token:token, Username:target}
    var reply cmd.FollowReply
    err, replyInf := ClientCall("Service." + "Follow", args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.FollowReply)))
}
