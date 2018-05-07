package cmd

type RegisterUserArgs struct {
    CommandId string
    Username string
    Password string
}

type RegisterUserReply struct {
    Ok bool
    Error string
    Token string
}

type RegisterUser struct {
    Args *RegisterUserArgs
    Channel chan *RegisterUserReply
}
