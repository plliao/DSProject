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
    err := client.InitOnce(network, address)
    for ; err != nil; {
        time.Sleep(1000 * time.Millisecond)
        //fmt.Print("Attempt to Connect with " + address + "\n")
        err = client.InitOnce(network, address)
    }
}

func (client *RaftClient) InitOnce(network string, address string) error {
    rpcClient, err := rpc.DialHTTPPath(network, address, rpc.DefaultRPCPath + address)
    client.rpcClient = rpcClient
    return err
}

func (client *RaftClient) AppendEntry(
        term int, leaderId int, prevLogIndex int, prevLogTerm int,
        command string, commandTerm int, commitIndex int) (AppendEntryReply, error) {
    args := AppendEntryArgs{
        Term:term,
        LeaderId:leaderId,
        PrevLogIndex:prevLogIndex,
        PrevLogTerm:prevLogTerm,
        Command:command,
        CommandTerm:commandTerm,
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
