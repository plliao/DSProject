package cmd

type UserLoginArgs struct {
    CommandId string
    Username string
    Password string
}

type UserLoginReply struct {
    Ok bool
    Error string
    Token string
}

type UserLogin struct {
    Args *UserLoginArgs
    Channel chan *UserLoginReply
}
