package backEnd

import (
    "fmt"
    "reflect"
    "time"
)

func (srv *Server) leaderInit() {
    srv.raft.isLeader = true
    for i, _ := range srv.addressBook {
        srv.nextIndexs[i] = len(srv.raft.logs)
    }
    srv.commitChan = make(chan int, 100)
    srv.commandLogs = make(map[string]reflect.Value)
    go srv.commitHandler()
    for i, _ := range srv.addressBook {
        if i != srv.id {
            go srv.followerHandler(i)
        }
    }
}

func (srv *Server) leaderShutDown() {
    close(srv.commitChan)
    srv.raft.isLeader = false
}

func (srv *Server) commitHandler() {
    indexCount := make(map[int]int)
    for index := range(srv.commitChan) {
        if srv.raft.commitIndex > index {
            continue
        }

        if _, ok := indexCount[index]; !ok {
            indexCount[index] = 1
        }
        indexCount[index] = indexCount[index] + 1

        if indexCount[index] > srv.getMajority() {
            for commitIndex := srv.raft.commitIndex + 1; commitIndex <= index; commitIndex++ {
                fmt.Printf("Exec %v\n", srv.commandLogs)
                encodedCmd := srv.raft.logs[commitIndex]
                results := srv.exec(encodedCmd)
                srv.raft.commitIndex = commitIndex

                commandId := srv.cmdFactory.GetCommandId(encodedCmd)
                if cmdValue, ok := srv.commandLogs[commandId]; ok {
                    srv.replyWithResults(cmdValue, results)
                    delete(srv.commandLogs, commandId)
                }

                delete(indexCount, commitIndex)
            }
        }
    }
}

func (srv *Server) followerHandler(index int) {
    fmt.Print("Start to Connect with " + srv.addressBook[index] + "\n")
    client := RaftClient{address:srv.addressBook[index]}
    client.Init(srv.network, srv.addressBook[index])
    fmt.Print("Successfully Connect with " + srv.addressBook[index] + "\n")
    delay := 1
    for {
        if !srv.raft.isLeader {
            break
        }
        nextIndex := srv.nextIndexs[index]
        var command string

        if srv.raft.index < nextIndex {
            command = ""
            nextIndex = srv.raft.index + 1
            time.Sleep(time.Duration(delay) * srv.timeout)
            delay++
        } else {
            command = srv.raft.logs[nextIndex]
        }

        _, commandTerm := srv.raft.getIndexAndTerm(nextIndex)
        preLogIndex, preLogTerm := srv.raft.getIndexAndTerm(nextIndex - 1)

        reply, err := client.AppendEntry(
            srv.raft.term,
            srv.id,
            preLogIndex,
            preLogTerm,
            command,
            commandTerm,
            srv.raft.commitIndex,
        )
        //fmt.Printf("Replicate to index:%v, term:%v, prevLogIndex:%v, preLogTerm:%v, command:%v, commit:%v\n", index, srv.raft.term, preLogIndex, preLogTerm, command, srv.raft.commitIndex)

        if err != nil {
            fmt.Print(err)
            client.Init(srv.network, srv.addressBook[index])
            continue
        }
        if reply.Term > srv.raft.term {
            srv.raft.toFollowerChan <- reply.Term
            break
        }
        if reply.Success {
            if command != "" {
                srv.nextIndexs[index]++
                if commandTerm == srv.raft.term {
                    srv.commitChan <- nextIndex
                }
            }
        } else {
            srv.nextIndexs[index]--
        }
    }
}
