package backEnd

import (
    "time"
    "sort"
    "html/template"
    "strings"
    "backEnd/cmd"
)

type User struct {
    Username string
    Password string
    token string
    Articles []*cmd.Article
    following map[string]*User
    followers map[string]*User
}

func (user *User) Init() {
    user.Articles = make([]*cmd.Article, 0)
    user.following = make(map[string]*User)
    user.following[user.Username] = user
    user.followers = make(map[string]*User)
    user.token = ""
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

func (user *User) Post(content string) {
    article := &cmd.Article{
        Content:content,
        Author:user.Username,
        Timestamp:time.Now(),
    }
    user.Articles = append(user.Articles, article)
}

func (follower *User) Follow(user *User) {
    follower.following[user.Username] = user
    user.followers[follower.Username] = follower
}

func (follower *User) UnFollow(user *User) {
    if _, ok := follower.following[user.Username]; ok {
        delete(follower.following, user.Username)
        delete(user.followers, follower.Username)
    }
}

func (follower *User) GetMyContent() []*cmd.Article {
    contents := make([]*cmd.Article, 0, 100)
    for _, user := range follower.following {
        contents = append(contents, user.Articles...)
    }
    sort.Slice(contents, func (i, j int) bool {
        return contents[i].Timestamp.After(contents[j].Timestamp)
    })
    return contents
}
