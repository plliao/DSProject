package cmd

type UserLoginArgs struct {
    Username string
    Password string
}

type UserLoginReply struct {
    Ok bool
    Token string
}

type UserLogin struct {
    Args *UserLoginArgs
    Channel chan *UserLoginReply
}
