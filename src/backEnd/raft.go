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

type RequestVoteArgs struct {
    term int
    candidateId int
    lastLogIndex int
    lastLogTerm int
}

type RequestVoteReply struct {
    term int
    voteGranted bool
}

func (raft *Raft) AppendEntry(args AppendEntryArgs, reply *AppendEntryReply) error {
    //TODO
    reply.success = false
    if args.term < raft.term || raft.logTerms[prevLogIndex] != prevLogTerm{
        return 
    }
    if len(raft.logs)-1 == args.prevLogIndex{
        raft.logs = append(raft.logs, command)
    }
    if commitIndex > raft.commitIndex {
        raft.commitIndex = min(commitIndex, len(logs)-1)
    }
    return nil
}

func (raft *Raft) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) error {
    //TODO
    return nil
}
