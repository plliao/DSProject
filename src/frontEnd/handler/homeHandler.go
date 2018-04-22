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

    client, dialerr := srv.ClientConnect()
    if(dialerr != nil){
        log.Fatal("LoginRPC:", dialerr)
    }
    user := &User{ Username:username , token:token}
    if(user.token != ""){
        if (post != "") {
            err, reply := ClientPostRPC(user.token, post, client)
            if(err != nil ){
                log.Fatal("PostRPC:", err)
            }
            if(!reply.Ok){
                http.Redirect(w, r, "/login/", http.StatusFound)
                return
            }
        }
        if (logout != "") {
            err, reply := ClientLogoutRPC(user.token, client)
            if(err == nil && reply.Ok){
                http.Redirect(w, r, "/login/", http.StatusFound)
                return
            }
            if(err != nil){
                log.Fatal("LoginRPC:", err)
            }
        }
    }else{
        var clientReply ClientReply
        var err error
        if(choose == "Sign up"){
            err, clientReply = ClientRegisterUserRPC(username, password, client)
        } else if(choose == "Log in"){
            err, clientReply = ClientUserLoginRPC(username, password, client)
        }
        if(err == nil && clientReply.Ok){
            user.token = clientReply.Token
        }else {
            if(err != nil){
                log.Fatal("Signup or Login RPC:", err)
            }
            loginURLValues := url.Values{}
            loginURLValues.Set("message", clientReply.Error)
            http.Redirect(w, r, "/login/?" + loginURLValues.Encode(), http.StatusFound)
            return
        } 
    }
    err, reply := ClientGetMyContentRPC(user.token, client)
    if(err == nil && reply.Ok){
        user.Articles = reply.Articles
    }
    if(err != nil){
        log.Fatal("GetMyContentRPC:", err)
    }
	server.RenderTemplate(w, srv, "home", user)
}

