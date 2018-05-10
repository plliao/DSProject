package handler

import (
	//"net/http"
    //"net/url"
    "backEnd/cmd"
    "net/rpc"
    "frontEnd/server"
    "strings"
    "crypto/rand"
    //"html/template"
    //"log"
    "fmt"
    "reflect"
    "time"
)

func ClientCall(service string, args interface{}, replyType reflect.Type, srv *server.Server) (error, interface{}){
    token := make([]byte, 6)
    rand.Read(token)
    ID := fmt.Sprintf("%x", token)

    stype := reflect.ValueOf(args)
    stype.Elem().FieldByName("CommandId").SetString(ID)

    NotLeader := "Not Leader: "

    for {
        address, network := srv.GetConnectInfo()
        fmt.Printf("\n%v: Dial to %v ", ID, address)
        client, errDial := rpc.DialHTTPPath(network, address, rpc.DefaultRPCPath + address)
        if errDial == nil {
            reply := reflect.New(replyType)
            errRPCChan := make(chan error, 1)
            var errRPC error
            go func(){
                errRPCChan <- client.Call(service, args, reply.Interface())
            }()
            select {
                case errRPC = <-errRPCChan:
                    fmt.Printf(", Get reply from %v", address)

                    ok := reply.Elem().Field(0).Interface().(bool)
                    message := reply.Elem().Field(1).Interface().(string)

                    if errRPC != nil {
                        fmt.Printf("\n%v\n", errRPC.Error())
                    } else if ok == false && strings.HasPrefix(message, NotLeader) {
                        address = message[len(NotLeader):]
                        if address != "" {
                            fmt.Printf(", Get new Address %v", address)
                            srv.SetConnectInfo(address, network)
                            continue
                        }
                    } else {
                        dupReply := reply.Interface()
                        return errRPC, dupReply
                    }
                case <-time.After(5000 * time.Millisecond):
                    fmt.Printf(", Timeout")
            }
        } else {
            fmt.Printf("\n%v\n", errDial.Error())
        }
        fmt.Printf(", Try next address")
        srv.TryNextAddress()
        time.Sleep(500 * time.Millisecond)
    }
}

func ClientPostRPC(token string, post string, srv *server.Server) (error, cmd.PostReply){
	args := cmd.PostArgs { Token:token, Content:post }
    var reply cmd.PostReply
    err, replyInf := ClientCall("Service." + "Post", &args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.PostReply)))
}

func ClientLogoutRPC(token string, srv *server.Server) (error, cmd.UserLogoutReply){
	args := cmd.UserLogoutArgs{ Token:token }
    var reply cmd.UserLogoutReply
	err, replyInf := ClientCall("Service." + "UserLogout", &args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.UserLogoutReply)))
}

func ClientRegisterUserRPC(username string, password string, srv *server.Server) (error, ClientReply){
	args := cmd.RegisterUserArgs{ Username:username, Password:password }
    var reply ClientReply
    err, replyInf := ClientCall("Service." + "RegisterUser", &args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*ClientReply))
}

func ClientUserLoginRPC(username string, password string, srv *server.Server) (error, ClientReply){
	args := cmd.UserLoginArgs{ Username:username, Password:password }
    var reply ClientReply
    err, replyInf := ClientCall("Service." + "UserLogin", &args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*ClientReply))
}

func ClientGetMyContentRPC(token string, srv *server.Server) (error, cmd.GetMyContentReply){
	args := cmd.GetMyContentArgs{ Token:token }
    var reply cmd.GetMyContentReply
    err, replyInf := ClientCall("Service." + "GetMyContent", &args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.GetMyContentReply)))
}

func ClientGetFollowerRPC(token string, srv *server.Server) (error, cmd.GetFollowerReply){
	args := cmd.GetFollowerArgs{ Token:token }
    var reply cmd.GetFollowerReply
    err, replyInf := ClientCall("Service." + "GetFollower", &args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.GetFollowerReply)))
}

func ClientDeleteUserRPC(token string, srv *server.Server) (error, cmd.DeleteUserReply){
	args := cmd.DeleteUserArgs{ Token:token }
    var reply cmd.DeleteUserReply
    err, replyInf := ClientCall("Service." + "DeleteUser", &args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.DeleteUserReply)))
}

func ClientUnFollowRPC(token string, target string, srv *server.Server) (error, cmd.UnFollowReply){
	args := cmd.UnFollowArgs{ Token:token, Username:target}
    var reply cmd.UnFollowReply
    err, replyInf := ClientCall("Service." + "UnFollow", &args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.UnFollowReply)))
}

func ClientFollowRPC(token string, target string, srv *server.Server) (error, cmd.FollowReply){
	args := cmd.FollowArgs{ Token:token, Username:target}
    var reply cmd.FollowReply
    err, replyInf := ClientCall("Service." + "Follow", &args, reflect.TypeOf(reply), srv)
    return err, *(replyInf.(*(cmd.FollowReply)))
}
