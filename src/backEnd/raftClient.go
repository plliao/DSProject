package backEnd

import(
    "net/rpc"
    "reflect"
    "time"
)

type RaftClient struct{
    address string
    rpcClient *rpc.Client
}

func (client *RaftClient) Init(network string, address string) {
    rpcClient, err := rpc.DialHTTP(network, address)
    for ; err != nil; {
        time.Sleep(1000)
        rpcClient, err = rpc.DialHTTP(network, address)
    }
    client.rpcClient = rpcClient
}

func (client *RaftClient) AppendEntry(
        term int, leaderId int, prevLogIndex int, prevLogTerm int,
        command reflect.Value, commitIndex int) (AppendEntryReply, error) {
    args := AppendEntryArgs{
        term:term,
        leaderId:leaderId,
        prevLogIndex:prevLogIndex,
        prevLogTerm:prevLogTerm,
        command:command,
        commitIndex:commitIndex,
    }
    reply := AppendEntryReply{}
    err := client.rpcClient.Call("Raft.AppendEntry", args, &reply)
    return reply, err
}
