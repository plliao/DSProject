package user

type User struct {
    Username string
    Password string
    Articles []string
    Following []*User
}
