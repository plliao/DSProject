package cmd

type UnFollowArgs struct {
    Token string
    Username string
}

type UnFollowReply struct {
    Ok bool
    Error string
}

type UnFollow struct {
    Args *UnFollowArgs
    Channel chan *UnFollowReply
}
