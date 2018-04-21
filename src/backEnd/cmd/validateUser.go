package cmd

import (
)

type ValidateUserArgs struct {
    Username string
    Password string
}

type ValidateUserReply struct {
    Ok bool
    Err error
}

type ValidateUser struct {
    Args *ValidateUserArgs
    Channel chan *ValidateUserReply
}
