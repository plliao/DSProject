package cmd

type UserLogoutArgs struct {
    CommandId string
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
