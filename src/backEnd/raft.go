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
    if args.Term < raft.term || raft.logTerms[args.PrevLogIndex] != args.PrevLogTerm{
        reply.Success = false
    }
    if len(raft.logs)-1 == args.PrevLogIndex{
        raft.logs = append(raft.logs, args.Command)
        reply.Success = true
    }else if len(raft.logs) >= args.PrevLogIndex+1 && raft.logs[args.PrevLogIndex+1] != args.Command{
        raft.logs = raft.logs[:args.PrevLogIndex]
        reply.Success = false
    }
    if args.CommitIndex > raft.commitIndex {
        for i:= raft.commitIndex; i<= args.CommitIndex; i++{
            //exec
        }
        if args.CommitIndex < len(raft.logs) - 1 {
            raft.commitIndex = args.CommitIndex
        } else {
            raft.commitIndex = len(raft.logs) - 1
        }
    }

    reply.Term = raft.term
    reply.Success = true
    return nil
}

func (raft *Raft) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) error {
    if args.Term < raft.term {
        reply.VoteGranted = false
    } else if ((raft.voteFor < 0 || raft.voteFor == args.CandidateId) &&
            len(raft.logs)-1 <= args.LastLogIndex &&
            raft.logTerms[len(raft.logTerms)-1] <= args.LastLogTerm) {
        reply.VoteGranted = true
        reply.Term = raft.term
        raft.voteFor = args.CandidateId
    }
    return nil
}
