package handler

/*import (
	"net/http"
    "net/url"
    "errors"
)*/

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

type SignupArg struct {
    Username string
    Password string
}

type LoginArg struct {
    Username string
    Password string
}

type SignupReply struct {
    ok bool
    Articles []Article
}

type LoginReply struct {
    ok bool
    Articles []Article
}

type PostArg struct {
    Username string
    Post string
}

type PostReply struct {
    ok bool
}

type GetContentArg struct {
    Username string
}

type GetContentReply struct {
    ok bool
    Articles []Article
}

type LogoutArg struct {
	Username string
}

type LogoutReply struct {
	ok bool
}

type DeleteArg struct {
	Username string
}

type DeleteReply struct {
	ok bool
}
