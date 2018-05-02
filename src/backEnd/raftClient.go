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

func (client *RaftClient) AppendEntry(term int, index int, commit int, command reflect.Value) {

}
