package backend

import(
    "net/rpc"
    "reflect"
)

type RaftClient struct{
    address string
    rpcClient *rpc.Client
}

func (client *RaftClient) Init() {
    rpcClient, err := rpc.DialHTTP("tcp", srv.addressBook[index])
    for ; err != nil; {
        time.Sleep(1000)
        rpcClient, err = rpc.DialHTTP("tcp", srv.addressBook[index])
    }
    client.rpcClient = rpcClient
}

func (client *RaftClient) AppendEntry(
        term int, leaderId int, prevLogIndex int, prevLogTerm int,
        command reflect.Value, commit int) AppendEntryReply {
    args := AppendEntryArgs{
        term:term,
        leaderId:leaderId,
        prevLogIndex:prevLogIndex,
        prevLogTerm:prevLogTerm,
        command:command,
        commit:commit,
    }
    reply = AppendEntryReply{}
    client.call("Raft.AppendEntry", args, &reply)
    return reply
}
