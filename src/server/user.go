package server

type User struct {
    Username string
    Password string
    Articles []string
    Following []*User
}
