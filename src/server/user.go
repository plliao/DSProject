package server

import (
    "time"
    "sort"
)

type Article struct {
    Content string
    Author string
    Timestamp time.Time
}

func (article *Article) GetTimeWithUnixDateFormat() string {
    return article.Timestamp.Format(time.UnixDate)
}

type User struct {
    Username string
    Password string
    Articles []*Article
    following map[string]*User
}

func (user *User) Init() {
    user.Articles = make([]*Article, 0)
    user.following = make(map[string]*User)
    user.following[user.Username] = user
}

func (user *User) Post(content string) {
    article := &Article{
        Content:content,
        Author:user.Username,
        Timestamp:time.Now(),
    }
    user.Articles = append(user.Articles, article)
}

func (follower *User) Follow(user *User) {
    follower.following[user.Username] = user
}

func (follower *User) UnFollow(user *User) {
    if _, ok := follower.following[user.Username]; ok {
        delete(follower.following, user.Username)
    }
}

func (follower *User) GetMyContent() []*Article {
    contents := make([]*Article, 0, 100)
    for _, user := range follower.following {
        contents = append(contents, user.Articles...)
    }
    sort.Slice(contents, func (i, j int) bool {
        return contents[i].Timestamp.After(contents[j].Timestamp)
    })
    return contents
}
