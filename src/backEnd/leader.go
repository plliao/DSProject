package backEnd

import (
    "fmt"
    "time"
    "reflect"
)

func (srv *Server) leaderInit() {
    srv.raft.leaderId = srv.id
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
    srv.appendCommand(reflect.ValueOf(srv.cmdFactory.MakeDummyCommand()))
    srv.raft.isLeader = true
}

func (srv *Server) leaderShutDown() {
    close(srv.commitChan)
    for _, cmdValue := range srv.commandLogs {
        srv.replyNotLeader(cmdValue)
    }
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
                encodedCmd := srv.raft.logs[commitIndex]
                cmdTerm := srv.raft.logTerms[commitIndex]
                log := fmt.Sprintf("CommitIndex %v, term %v, leader term %v: %v\n", commitIndex, cmdTerm, srv.raft.term, encodedCmd)
                srv.logger.WriteString(log)
                srv.logger.Flush()
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
    client := RaftClient{address:srv.addressBook[index]}
    client.Init(srv.network, srv.addressBook[index])
    delay := 1
    for {
        nextIndex := srv.nextIndexs[index]
        var command string

        if srv.raft.index < nextIndex {
            command = ""
            srv.nextIndexs[index] = srv.raft.index + 1
            nextIndex = srv.nextIndexs[index]
            time.Sleep(time.Duration(delay) * srv.timeout)
            //delay++
            fmt.Printf("%v", len(srv.raft.logs))
        } else {
            command = srv.raft.logs[nextIndex]
        }

        _, commandTerm := srv.raft.getIndexAndTerm(nextIndex)
        preLogIndex, preLogTerm := srv.raft.getIndexAndTerm(nextIndex - 1)

        srv.rwLock.RLock()
        if !srv.raft.isLeader {
            srv.rwLock.RUnlock()
            break
        }
        reply, err := client.AppendEntry(
            srv.raft.term,
            srv.id,
            preLogIndex,
            preLogTerm,
            command,
            commandTerm,
            srv.raft.commitIndex,
        )
        srv.rwLock.RUnlock()

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
            if srv.nextIndexs[index] < 0 {
                srv.nextIndexs[index] = 0
            }
        }
    }
}
