package backEnd

import (
    "reflect"
)

type Raft struct {
    isLeader bool

    term int
    voteFor int
    commitIndex int
    index int

    logs []reflect.Value
    logTerms []int
}

type AppendEntryArgs struct {
    term int
    leaderId int
    prevLogIndex int
    prevLogTerm int
    command reflect.Value
    commitIndex int
}

type AppendEntryReply struct {
    term int
    success bool
}
