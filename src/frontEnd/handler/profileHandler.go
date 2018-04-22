package handler

import (
	"net/http"
    "net/url"
    "frontEnd/server"
    "backEnd/cmd"
    //"errors"
    //"strings"
    //"html/template"
    //"log"
)

/*type FollowButton struct {
    Name string
    Action string
    User *User
}*/

func ProfileHandler(w http.ResponseWriter, r *http.Request, srv *server.Server){
    token := r.FormValue("Auth")
    username := r.FormValue("name")
    user := &User{ token:token , Username:username}

    args := cmd.GetFollowerArgs{ Token:user.token }
    var reply cmd.GetFollowerReply 
    srv.SrvClient.Call("Service." + "GetFollower", args, &reply)
    errmsg := ""
    userMap := make(map[string]int)

    if(!reply.Ok){
        errmsg = "Wrong User."
        loginURLValues := url.Values{}
        loginURLValues.Set("message", errmsg)
        http.Redirect(w, r, "/login/?" + loginURLValues.Encode(), http.StatusFound)
    }
    for i, u := range(reply.Relationships){
        tmpAction := "Follow"
        if(u.Following){
            tmpAction = "Unfollow"
        }
        user.Others = append(user.Others, Relationship{Username:u.Username, Action:tmpAction})
        userMap[u.Username] = i
    }

    action := r.FormValue("FollowOrNot")
    target := r.FormValue("Target")
    deleteAccount := r.FormValue("delete")

    if deleteAccount != "" {
    	args := cmd.DeleteUserArgs{ Token:token }
        var reply cmd.DeleteUserReply
        err := srv.SrvClient.Call("Service." + "DeleteUser", args, &reply)
        errmsg := ""
        if(err != nil || !reply.Ok){
            errmsg = "Delete Failed. Please log in again."
        }
        loginURLValues := url.Values{}
        loginURLValues.Set("message", errmsg)
        http.Redirect(w, r, "/login/?" + loginURLValues.Encode(), http.StatusFound)
        return
    }
    if action == "Unfollow" {
        args := cmd.UnFollowArgs{ Token:token, Username:target}
        var reply cmd.UnFollowReply
        srv.SrvClient.Call("Service." + "UnFollow", args, &reply)
        if(reply.Ok){
            user.Others[userMap[target]].Action = "Follow"
        }
    }else if(action == "Follow"){
        args := cmd.FollowArgs{ Token:token, Username:target}
        var reply cmd.FollowReply
        srv.SrvClient.Call("Service." + "Follow", args, &reply)
        if(reply.Ok){
            user.Others[userMap[target]].Action = "Unfollow"
        }
    }
    server.RenderTemplate(w, srv, "profile", user)
}

