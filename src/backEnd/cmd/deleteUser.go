package cmd

type DeleteUserArgs struct {
    Token string
}

type DeleteUserReply struct {
    Ok bool
    Error string
}

type DeleteUser struct {
    Args *DeleteUserArgs
    Channel chan *DeleteUserReply
}
