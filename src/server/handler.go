package server

import (
	"net/http"
    "net/url"
    "errors"
)

type ServerHandlerFunc func(http.ResponseWriter, *http.Request, *Server)

type HandlerFuncFactory struct {
}

func (factory *HandlerFuncFactory) CreateByServerHandlerFunc(
        serverHandler ServerHandlerFunc, srv *Server) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        serverHandler(w, r, srv)
    }
}

type LoginPage struct {
	Message string
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

func LoginHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
    message := r.URL.Query().Get("message")
	  RenderTemplate(w, srv, "login", LoginPage{message})
}

func HomeHandler(w http.ResponseWriter, r *http.Request, srv *Server) {
    token := r.FormValue("Auth")
    var user *User
    if srv.ValidateAuth(token) {
        user = srv.tokens[token]
        post := r.FormValue("article")
        if post != "" {
            user.Post(post)
        }
        logout := r.FormValue("logout")
        if logout != "" {
            srv.UserLogout(user)
            http.Redirect(w, r, "/login/", http.StatusFound)
            return
        }
    } else {
        username := r.FormValue("name")
        password := r.FormValue("password")
        choose := r.FormValue("choose")

        ok := false
        var err error
        if(choose == "Log in"){
            ok, err = srv.ValidateUser(username, password)
        } else if (choose == "Sign up") {
            ok, err = srv.RegisterUser(username, password)
        } else {
            err = errors.New("")
        }

        if !ok {
            loginURLValues := url.Values{}
            loginURLValues.Set("message", err.Error())
            http.Redirect(w, r, "/login/?" + loginURLValues.Encode(), http.StatusFound)
            return
        }
        user = srv.users[username]
        srv.UserLogin(user)
    }
	RenderTemplate(w, srv, "home", user)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request, srv *Server){
    token := r.FormValue("Auth")
    var user *User
    var buttons []FollowButton
    var profile ProfilePage
    if srv.ValidateAuth(token) {
        user = srv.tokens[token]
        action := r.FormValue("FollowOrNot")
        target := r.FormValue("Target")
        deleteAccount := r.FormValue("delete")

        if deleteAccount != "" {
            srv.DeleteUser(user)
            http.Redirect(w, r, "/login/", http.StatusFound)
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
    } else {
        http.Redirect(w, r, "/login/", http.StatusFound)
        return
    }
    RenderTemplate(w, srv, "profile", profile)
}

