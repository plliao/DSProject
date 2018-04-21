package cmd

type PostArgs struct {
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
