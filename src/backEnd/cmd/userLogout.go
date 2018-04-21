package cmd

type UserLogoutArgs struct {
    Token string
}

type UserLogoutReply struct {
    Ok bool
    Error string
}

type UserLogout struct {
    Args *UserLogoutArgs
    Channel chan *UserLogoutReply
}
