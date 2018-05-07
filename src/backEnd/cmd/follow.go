package cmd

type FollowArgs struct {
    CommandId string
    Token string
    Username string
}

type FollowReply struct {
    Ok bool
    Error string
}

type Follow struct {
    Args *FollowArgs
    Channel chan *FollowReply
}
