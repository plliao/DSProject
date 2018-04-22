package handler

import (
	"net/http"
    "net/url"
    "frontEnd/server"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request, srv *server.Server){
    token := r.FormValue("Auth")
    username := r.FormValue("name")
    user := &User{ token:token , Username:username}
    errmsg, reply := ClientGetFollowerRPC(token, srv)
    if(!reply.Ok){
        loginURLValues := url.Values{}
        loginURLValues.Set("message", errmsg)
        http.Redirect(w, r, "/login/?" + loginURLValues.Encode(), http.StatusFound)
        return
    }
    userMap := make(map[string]int)
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
    	errmsg, reply := ClientDeleteUserRPC(user.token, srv)
        if(reply.Ok){
            loginURLValues := url.Values{}
            loginURLValues.Set("message", errmsg)
            http.Redirect(w, r, "/login/?" + loginURLValues.Encode(), http.StatusFound)
            return
        }
    }
    if action == "Unfollow" {
        err, reply := ClientUnFollowRPC(user.token, target, srv)
        if(err == nil && reply.Ok){
            user.Others[userMap[target]].Action = "Follow"
        }
    }else if(action == "Follow"){
        err, reply := ClientFollowRPC(user.token, target, srv)
        if(err == nil && reply.Ok){
            user.Others[userMap[target]].Action = "Unfollow"
        }
    }
    server.RenderTemplate(w, srv, "profile", user)
}

