package cmd

type RegisterUserArgs struct {
    Username string
    Password string
}

type RegisterUserReply struct {
    Ok bool
    Token string
    Err error
}

type RegisterUser struct {
    Args *RegisterUserArgs
    Channel chan *RegisterUserReply
}
