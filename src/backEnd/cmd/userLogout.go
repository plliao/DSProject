package cmd

type UserLogoutArgs struct {
    Token string
}

type UserLogoutReply struct {
    Ok bool
}

type UserLogout struct {
    Args *UserLogoutArgs
    Channel chan *UserLogoutReply
}
