package backend

import (
    "reflect"
)

type Raft struct {
    term int
    voteFor int
    commitIndex int
    index int

    logs []reflect.Value
}

type AppendEntryArgs {
    term int
    leaderId int
    prevLogIndex int
    prevLogTerm int
    command reflect.Value
    commit int
}

type AppendEntryReply {
    term int
    success bool
}
