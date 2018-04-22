package cmd

type Relationship struct {
    Username string
    Following bool
}

type GetFollowerArgs struct {
    Token string
}

type GetFollowerReply struct {
    Ok bool
    Error string
    Relationships []*Relationship
}

type GetFollower struct {
    Args *GetFollowerArgs
    Channel chan *GetFollowerReply
}
