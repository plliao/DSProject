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
    Token string
}

type GetMyContentReply struct {
    Articles []*Article
    Ok bool
}

type GetMyContent struct {
    Args *GetMyContentArgs
    Channel chan *GetMyContentReply
}
