package backEnd

import(
    "net/rpc"
    "time"
    "fmt"
)

type RaftClient struct{
    address string
    rpcClient *rpc.Client
}

func (client *RaftClient) Init(network string, address string) {
    fmt.Print("Start to Connect with " + address + "\n")
    rpcClient, err := rpc.DialHTTP(network, address)
    for ; err != nil; {
        time.Sleep(1000 * time.Millisecond)
        fmt.Print("Attempt to Connect with " + address + "\n")
        rpcClient, err = rpc.DialHTTP(network, address)
    }
    client.rpcClient = rpcClient
}

func (client *RaftClient) AppendEntry(
        term int, leaderId int, prevLogIndex int, prevLogTerm int,
        command string, commitIndex int) (AppendEntryReply, error) {
    args := AppendEntryArgs{
        Term:term,
        LeaderId:leaderId,
        PrevLogIndex:prevLogIndex,
        PrevLogTerm:prevLogTerm,
        Command:command,
        CommitIndex:commitIndex,
    }
    reply := AppendEntryReply{}
    err := client.rpcClient.Call("Raft.AppendEntry", args, &reply)
    return reply, err
}

func (client *RaftClient) RequestVote(
        term int, candidateId int, lastLogIndex int, lastLogTerm int) (RequestVoteReply, error) {
    args := RequestVoteArgs{
        Term:term,
        CandidateId:candidateId,
        LastLogIndex:lastLogIndex,
        LastLogTerm:lastLogTerm,
    }
    reply := RequestVoteReply{}
    err := client.rpcClient.Call("Raft.RequestVote", args, &reply)
    return reply, err
}
