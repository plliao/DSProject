package cmd

import (
    "time"
)

type Article struct {
    Content string
    Author string
    Timestamp time.Time
}

func (article *Article) GetTimeWithUnixDateFormat() string {
    return article.Timestamp.Format(time.UnixDate)
}

type GetMyContentArgs struct {
    CommandId string
    Token string
}

type GetMyContentReply struct {
    Ok bool
    Error string
    Articles []*Article
}

type GetMyContent struct {
    Args *GetMyContentArgs
    Channel chan *GetMyContentReply
}
