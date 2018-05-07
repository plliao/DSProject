package cmd

type PostArgs struct {
    CommandId string
    Token string
    Content string
}

type PostReply struct {
    Ok bool
    Error string
}

type Post struct {
    Args *PostArgs
    Channel chan *PostReply
}
