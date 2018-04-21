package handler

import (
	"net/http"
    "time"
    "frontEnd/server"
    "backEnd/cmd"
    //"errors"
    "strings"
    "html/template"
)

type User struct {
    Username string
    Password string
    Articles []Article
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

type Article struct {
    Content string
    Author string
    Timestamp time.Time
}

type ClientReply struct {
    Ok bool
    Token string
    Error string
}

func (article *Article) GetTimeWithUnixDateFormat() string {
    return article.Timestamp.Format(time.UnixDate)
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

    user := User{ Username:username }

    if(choose == "Sign up" || choose == "Log in"){
        var clientReply ClientReply
        var err error
        if(choose == "Sign up"){
            args := cmd.RegisterUserArgs{ Username:username, Password:password }
            var reply cmd.RegisterUserReply
            err = srv.SrvClient.Call("Server.RegisterUser", args, &reply)
            clientReply = ClientReply{ reply.Ok, reply.Token, reply.Error }
        } else if(choose == "Log in"){
            args := cmd.UserLoginArgs{ Username:username, Password:password }
            var reply cmd.UserLoginReply
            err = srv.SrvClient.Call("Server.ValidateUser", args, &reply)
            clientReply = ClientReply{ reply.Ok, reply.Token, reply.Error }
        }
        if(err == nil && clientReply.Ok){
            user.token = clientReply.Token
        }else {
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
            if(err != nil || !reply.Ok){
                http.Redirect(w, r, "/login/", http.StatusFound)
                return
            }
        }
        if (logout != "") {
            args := cmd.UserLogoutArgs{ Token:user.token }
            var reply cmd.UserLogoutReply
            err := srv.SrvClient.Call("Server.UserLogout", args, &reply)
            if(err == nil && reply.Ok){
                http.Redirect(w, r, "/login/", http.StatusFound)
            }
            return
        }
        args := cmd.GetMyContentArgs{ Token:user.token }
        var reply cmd.GetMyContentReply
        err := srv.SrvClient.Call("User.GetMyContent", args, &reply)
        if(err == nil && reply.ok){
            user.Articles = reply.Articles
        }
    }
	server.RenderTemplate(w, srv, "home", user)
}

