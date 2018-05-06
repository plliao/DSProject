package backEnd

import (
    "fmt"
    "time"
)

type Raft struct {
    isLeader bool
    leaderId int

    term int
    voteFor int
    commitIndex int
    index int

    logs []string
    logTerms []int

    toExecChan chan int
    heartBeatChan chan time.Time
    toFollowerChan chan int
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

func (raft *Raft) appendCommand(command string, term int) {
    raft.logs = append(raft.logs, command)
    raft.logTerms = append(raft.logTerms, term)
    raft.index = len(raft.logs) - 1
}

func (raft *Raft) resetCommand(prevLogIndex int) {
    raft.logs = raft.logs[:prevLogIndex]
    raft.logTerms = raft.logTerms[:prevLogIndex]
    raft.index = len(raft.logs) - 1
}

func (raft *Raft) getLastIndexAndTerm() (int, int) {
    lastLogIndex := raft.index
    lastLogTerm := -1
    if lastLogIndex >= 0 {
        lastLogTerm = raft.logTerms[lastLogIndex]
    }
    return lastLogIndex, lastLogTerm
}

func (raft *Raft) AppendEntry(args AppendEntryArgs, reply *AppendEntryReply) error {
    raft.heartBeatChan <- time.Now()
    reply.Term = raft.term

    if args.Term < raft.term || args.PrevLogIndex >= len(raft.logs) ||
            args.PrevLogIndex > 0 && raft.logTerms[args.PrevLogIndex] != args.PrevLogTerm {
        reply.Success = false
        return nil
    }

    if raft.term < args.Term {
        raft.voteFor = -1
        raft.leaderId = args.LeaderId
        raft.toFollowerChan <- args.Term
    }

    if len(raft.logs) - 1 == args.PrevLogIndex{
        if args.Command != "" {
            fmt.Print("Append command " + args.Command + "\n")
            raft.appendCommand(args.Command, args.Term)
        } else {
            fmt.Print("HeartBeat\n")
        }
    } else if len(raft.logs) > args.PrevLogIndex + 1 {
        if raft.logs[args.PrevLogIndex + 1] != args.Command {
            raft.resetCommand(args.PrevLogIndex)
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
    reply.Success = true
    return nil
}

func (raft *Raft) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) error {
    fmt.Printf("Current term %v, voteFor %v\n", raft.term, raft.voteFor)
    fmt.Printf("Vote Request %v\n", args)
    lastLogIndex, lastLogTerm := raft.getLastIndexAndTerm()
    if args.Term < raft.term {
        reply.VoteGranted = false
        return nil
    }

    if args.Term > raft.term {
        raft.voteFor = -1
        raft.toFollowerChan <- args.Term
    }

    if (raft.voteFor < 0 || raft.voteFor == args.CandidateId) &&
            lastLogIndex <= args.LastLogIndex && lastLogTerm <= args.LastLogTerm {
        reply.VoteGranted = true
        reply.Term = raft.term
        raft.voteFor = args.CandidateId
    }
    return nil
}
