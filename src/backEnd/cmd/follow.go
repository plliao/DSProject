package cmd

type FollowArgs struct {
    Token string
    Username string
}

type FollowReply struct {
    Ok bool
}

type Follow struct {
    Args *FollowArgs
    Channel chan *FollowReply
}
