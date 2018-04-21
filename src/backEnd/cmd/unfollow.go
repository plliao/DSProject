package cmd

type UnFollowArgs struct {
    Token string
    Username string
}

type UnFollowReply struct {
    Ok bool
}

type UnFollow struct {
    Args *UnFollowArgs
    Channel chan *UnFollowReply
}
