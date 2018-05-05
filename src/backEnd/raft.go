package backEnd

import (
    "fmt"
    "time"
)

type Raft struct {
    isLeader bool

    term int
    voteFor int
    commitIndex int
    index int

    logs []string
    logTerms []int

    toExecChan chan int
    heartBeatChan chan time.Time
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
    raft.heartBeatChan <- time.Now()
    reply.Term = raft.term

    if args.Term < raft.term || args.PrevLogIndex >= len(raft.logs) ||
            args.PrevLogIndex > 0 && raft.logTerms[args.PrevLogIndex] != args.PrevLogTerm {
        reply.Success = false
        return nil
    }

    if len(raft.logs) - 1 == args.PrevLogIndex{
        if args.Command != "" {
            fmt.Print("Append command " + args.Command + "\n")
            raft.logs = append(raft.logs, args.Command)
            raft.logTerms = append(raft.logTerms, args.Term)
        } else {
            fmt.Print("HeartBeat\n")
        }
    } else if len(raft.logs) > args.PrevLogIndex + 1 {
        if raft.logs[args.PrevLogIndex + 1] != args.Command {
            raft.logs = raft.logs[:args.PrevLogIndex]
            raft.logTerms = raft.logTerms[:args.PrevLogIndex]
            reply.Success = false
            return nil
        }
    }

    if args.CommitIndex > raft.commitIndex {
        newCommitIndex := args.CommitIndex
        if newCommitIndex > len(raft.logs) - 1 {
            newCommitIndex = len(raft.logs) - 1
        }
        for i:= raft.commitIndex + 1; i<= newCommitIndex; i++{
            raft.toExecChan <- i
        }
        raft.commitIndex = newCommitIndex
    }
    raft.voteFor = -1
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
