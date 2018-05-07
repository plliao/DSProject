package cmd

type DeleteUserArgs struct {
    CommandId string
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
