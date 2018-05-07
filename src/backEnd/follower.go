package backEnd

import (
    "time"
    mrand "math/rand"
    "fmt"
)

func (srv *Server) followerInit() {
    srv.heartBeatChan = make(chan time.Time, 100)
    srv.toExecChan = make(chan int, 100)
    srv.raft.heartBeatChan = srv.heartBeatChan
    srv.raft.toExecChan = srv.toExecChan
    go srv.heartBeatHandler()
    go srv.execHandler()
}

func (srv *Server) followerShutDown() {
    srv.raft.toExecChan = nil
    srv.raft.heartBeatChan = nil
    close(srv.toExecChan)
    close(srv.heartBeatChan)
}

func (srv *Server) updateLastBeat() {
    for beatTime := range srv.raft.heartBeatChan {
        srv.lastBeatTime = beatTime
    }
}

func (srv *Server) startVote() bool {
    count := 1
    srv.raft.term = srv.raft.term + 1
    srv.raft.voteFor = srv.id
    countChan := make(chan int, len(srv.addressBook))
    for index, _ := range srv.addressBook {
        if srv.id == index {
            continue
        }
        go func(index int) {
            client := RaftClient{}
            err := client.InitOnce(srv.network, srv.addressBook[index])
            if err != nil {
                countChan <- 0
                return
            }
            lastLogIndex, lastLogTerm := srv.raft.getLastIndexAndTerm()
            reply, err := client.RequestVote(
                srv.raft.term,
                srv.id,
                lastLogIndex,
                lastLogTerm)
            if err == nil && reply.VoteGranted {
                countChan <- 1
            } else {
                countChan <- 0
            }
        }(index)
    }
    times := 1
    for result := range countChan {
        times++
        count += result
        if count > srv.getMajority() || times == len(srv.addressBook) {
            break
        }
    }
    if count > srv.getMajority() {
        return true
    }
    return false
}

func (srv *Server) heartBeatHandler(){
    go srv.updateLastBeat()
    for {
        time.Sleep(srv.timeout)
        randomTimeout := time.Duration(mrand.Intn(3) + 2) * srv.timeout
        if time.Now().Sub(srv.lastBeatTime) > randomTimeout {
            srv.lastBeatTime = time.Now()
            fmt.Print("Leader timeout\n")
            electionTimer := 10 * srv.timeout
            startVoteChan := make(chan bool, 1)
            go func(){
                startVoteChan <- srv.startVote()
            }()
            select {
                case voteRes := <-startVoteChan:
                    fmt.Printf("Election result: %v\n", voteRes)
                    if voteRes {
                        srv.followerShutDown()
                        srv.leaderInit()
                        return
                    }
                case <-time.After(electionTimer):
                    fmt.Println("election timeout")
            }
        }
    }
}

func (srv *Server) execHandler() {
    for execID :=  range srv.raft.toExecChan {
        srv.execCommit(execID)
    }
}
