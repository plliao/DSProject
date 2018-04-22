package handler

import (
	"net/http"
    "frontEnd/server"
    "backEnd/cmd"
)

type FollowButton struct {
    Name string
    Action string
    User *User
}

type ProfilePage struct {
    User *User
    Auth string
    FollowList []FollowButton
}

func ProfileHandler(w http.ResponseWriter, r *http.Request, srv *server.Server){
    var buttons []FollowButton
    token := r.FormValue("Auth")
    user := &User{ token:token }

    args := cmd.GetFollowerArgs{ Token:user.token }
    var reply cmd.GetFollowerReply 
    err := srv.SrvClient.Call("Service." + "GetFollower", args, &reply)
    errmsg := ""

    if(!reply.Ok){
        errmsg = "Wrong User."
        loginURLValues := url.Values{}
        loginURLValues.Set("message", errmsg)
        http.Redirect(w, r, "/login/?" + loginURLValues.Encode(), http.StatusFound)
    }
    user.following = reply.following
    user.unfollowing = reply.unfollowing

    var profile ProfilePage
    action := r.FormValue("FollowOrNot")
    target := r.FormValue("Target")
    deleteAccount := r.FormValue("delete")

    if deleteAccount != "" {
    	args := cmd.DeleteUserArgs{ Token:token }
        var reply DeleteReply
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
        user.UnFollow(srv.users[target])
    }else if(action == "Follow"){
        user.Follow(srv.users[target])
    }

    for u := range(srv.users){
        _, ok := user.following[u]
        if u != user.Username {
            if ok {
                buttons = append(buttons, FollowButton{Name:u, Action:"Unfollow", User:user})
            } else {
                buttons = append(buttons, FollowButton{Name:u, Action:"Follow", User:user})
            }
        }
    }
    profile = ProfilePage{User:user, FollowList:buttons}
    server.RenderTemplate(w, srv, "profile", profile)
}

