package handler

import (
	"net/http"
    "net/url"
    "frontEnd/server"
    "backEnd/cmd"
    //"errors"
    "strings"
    "html/template"
    "log"
)

type User struct {
    Username string
    Password string
    Articles []*cmd.Article
    token string
    //following map[string]*User
    //followers map[string]*User
}

func (user *User) Auth() template.HTML {
    htmlTokens := []string{
        "<input",
        "type=\"hidden\"",
        "name=\"Auth\"",
        "value=\"" + user.token + "\"",
        ">",
        "</input>",
    }
    return template.HTML(strings.Join(htmlTokens, " "))
}

type ClientReply struct {
    Ok bool
    Token string
    Error string
}

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

type LoginPage struct {
    Message string
}

func LoginHandler(w http.ResponseWriter, r *http.Request, srv *server.Server) {
    message := r.URL.Query().Get("message")
	server.RenderTemplate(w, srv, "login", LoginPage{message})
}

func HomeHandler(w http.ResponseWriter, r *http.Request, srv *server.Server) {
    username := r.FormValue("name")
    password := r.FormValue("password")
    choose := r.FormValue("choose")

    userref := User{ Username:username }
    user := &userref
    user.token = r.FormValue("Auth")

    if(choose == "Sign up" || choose == "Log in"){
        var clientReply ClientReply
        var err error
        if(choose == "Sign up"){
            args := cmd.RegisterUserArgs{ Username:username, Password:password }
            var reply cmd.RegisterUserReply
            err = srv.SrvClient.Call("Service." + "RegisterUser", args, &reply)
            clientReply = ClientReply{ reply.Ok, reply.Token, reply.Error }
        } else if(choose == "Log in"){
            args := cmd.UserLoginArgs{ Username:username, Password:password }
            var reply cmd.UserLoginReply
            err = srv.SrvClient.Call("Service." + "UserLogin", args, &reply)
            clientReply = ClientReply{ reply.Ok, reply.Token, reply.Error }
        }
        if(err == nil && clientReply.Ok){
            user.token = clientReply.Token
        }else {
            if(err != nil){
                log.Fatal("dialing:", err)
            }
            loginURLValues := url.Values{}
            loginURLValues.Set("message", clientReply.Error)
            http.Redirect(w, r, "/login/?" + loginURLValues.Encode(), http.StatusFound)
            return
        } 
    }else {
        post := r.FormValue("article")
        logout := r.FormValue("logout")
        if (post != "") {
            args := cmd.PostArgs { Token:user.token, Content:post }
            var reply cmd.PostReply
            err := srv.SrvClient.Call("User.Post", args, &reply)
            if(err != nil ){
                http.Redirect(w, r, "/login/", http.StatusFound)
                return
            }
        }
        if (logout != "") {
            args := cmd.UserLogoutArgs{ Token:user.token }
            var reply cmd.UserLogoutReply
            err := srv.SrvClient.Call("Service" + "UserLogout", args, &reply)
            if(err == nil && reply.Ok){
                http.Redirect(w, r, "/login/", http.StatusFound)
            }
            return
        }
        args := cmd.GetMyContentArgs{ Token:user.token }
        var reply cmd.GetMyContentReply
        err := srv.SrvClient.Call("Service" + "GetMyContent", args, &reply)
        if(err == nil && reply.Ok){
            user.Articles = reply.Articles
        }
    }
	server.RenderTemplate(w, srv, "home", user)
}

