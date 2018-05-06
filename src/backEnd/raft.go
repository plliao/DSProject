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
    CommandTerm int
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

func (raft *Raft) contains(index int, term int) bool {
    if index < 0 {
        return true
    }
    return raft.index >= index && raft.logTerms[index] == term
}

func (raft *Raft) match(index int, term int, command string) bool {
    if !raft.contains(index, term) {
        return false
    }
    return raft.logs[index] == command
}

func (raft *Raft) getIndexAndTerm(index int) (int, int) {
    if index < 0 {
        return -1, -1
    }
    if index > raft.index {
        return index, raft.term
    }
    return index, raft.logTerms[index]
}

func (raft *Raft) getLastIndexAndTerm() (int, int) {
    return raft.getIndexAndTerm(raft.index)
}

func (raft *Raft) AppendEntry(args AppendEntryArgs, reply *AppendEntryReply) error {
    raft.heartBeatChan <- time.Now()
    reply.Term = raft.term

    if args.Term < raft.term || !raft.contains(args.PrevLogIndex, args.PrevLogTerm) {
        reply.Success = false
        return nil
    }

    if raft.term < args.Term {
        raft.voteFor = -1
        raft.leaderId = args.LeaderId
        raft.toFollowerChan <- args.Term
    }

    currentIndex := args.PrevLogIndex + 1

    if len(raft.logs) == currentIndex {
        if args.Command != "" {
            fmt.Print("Append command " + args.Command + "\n")
            raft.appendCommand(args.Command, args.CommandTerm)
        } else {
            fmt.Print("HeartBeat\n")
        }
    } else {
        if !raft.match(currentIndex, args.CommandTerm, args.Command) {
            raft.resetCommand(args.PrevLogIndex)
            reply.Success = false
            return nil
        }
    }

    if args.CommitIndex > raft.commitIndex {
        newCommitIndex := args.CommitIndex
        if newCommitIndex > raft.index {
            newCommitIndex = raft.index
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
    reply.Term = raft.term
    if args.Term < raft.term {
        reply.VoteGranted = false
        return nil
    }

    lastLogIndex, lastLogTerm := raft.getLastIndexAndTerm()
    if args.Term > raft.term {
        raft.voteFor = -1
        raft.toFollowerChan <- args.Term
    }

    if (raft.voteFor < 0 || raft.voteFor == args.CandidateId) &&
            lastLogIndex <= args.LastLogIndex && lastLogTerm <= args.LastLogTerm {
        reply.VoteGranted = true
        raft.voteFor = args.CandidateId
    }
    return nil
}
