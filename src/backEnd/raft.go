package backEnd

import (
    "fmt"
)

type Raft struct {
    isLeader bool

    term int
    voteFor int
    commitIndex int
    index int

    logs []string
    logTerms []int
}

type AppendEntryArgs struct {
    Term int
    LeaderId int
    PrevLogIndex int
    PrevLogTerm int
    Command string
    CommitIndex int
}

type AppendEntryReply struct {
    Term int
    Success bool
}

type RequestVoteArgs struct {
    Term int
    CandidateId int
    LastLogIndex int
    LastLogTerm int
}

type RequestVoteReply struct {
    Term int
    VoteGranted bool
}

func (raft *Raft) AppendEntry(args AppendEntryArgs, reply *AppendEntryReply) error {
    fmt.Print("Receive AppendEntry\n")
    if args.term < raft.term || raft.logTerms[prevLogIndex] != prevLogTerm{
        reply.success = false
    }
    if len(raft.logs)-1 == args.prevLogIndex{
        raft.logs = append(raft.logs, command)
        reply.success = true
    }else if len(raft.logs) >= args.prevLogIndex+1 && raft.logs[prevLogIndex+1] != command{
        raft.logs = raft.logs[:prevLogIndex]
        reply.success = false
    }
    if args.commitIndex > raft.commitIndex {
        for i:= raft.commitIndex; i<= args.commitIndex; i++{
            //exec
        }
        raft.commitIndex = min(commitIndex, len(logs)-1)
    }

    reply.term = raft.term
    return nil
}

func (raft *Raft) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) error {
    if args.term < raft.term {
        reply.voteGranted = false
    }else if ((raft.voteFor < 0 || raft.voteFor == args.candidateId)
        && len(raft.logs)-1 <= args.lastLogIndex
        && raft.logTerms[len(raft.logTerms)-1] <= args.lastLogTerm){
        reply.voteGranted = true
        reply.term = raft.term
        raft.voteFor = args.candidateId
    }
    return nil
}
