package cmd

type DummyArgs struct {
    CommandId string
}

type DummyReply struct {
    Ok bool
    Error string
}

type Dummy struct {
    Args *DummyArgs
    Channel chan *DummyReply
}
