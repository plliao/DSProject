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

type Relationship struct {
    Username string
    Action string
}

type User struct {
    Username string
    Password string
    Articles []*cmd.Article
    token string
    Others []Relationship 
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
    token := r.FormValue("Auth")
    post := r.FormValue("article")
    logout := r.FormValue("logout")

    user := &User{ Username:username , token:token}
    if(token != ""){
        if (post != "") {
            err, reply := ClientPostRPC(token, post, srv)
            if(err != nil || !reply.Ok){
                http.Redirect(w, r, "/login/", http.StatusFound)
                return
            }
        }
        if (logout != "") {
            err, reply := ClientLogoutRPC(token, srv)
            if(err == nil && reply.Ok){
                http.Redirect(w, r, "/login/", http.StatusFound)
                return
            }
        }
    }else{
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
    }
    args := cmd.GetMyContentArgs{ Token:user.token }
    var reply cmd.GetMyContentReply
    err := srv.SrvClient.Call("Service." + "GetMyContent", args, &reply)
    if(err == nil && reply.Ok){
        user.Articles = reply.Articles
    }
	server.RenderTemplate(w, srv, "home", user)
}

