package user

type User struct {
    username string,
    password []byte,
    articles [][]byte,
    following []*User
}
