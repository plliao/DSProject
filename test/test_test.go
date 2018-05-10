package main

import(
    //"net/rpc"
    "time"
    "fmt"
    "frontEnd/server"
    "frontEnd/handler"
    "testing"
    "backEnd"
    //"backEnd/cmd"
    "strings"
    "bufio"
    "os"
    "log"
    "sync"
)

type Config struct {
    mu          sync.Mutex
    backEndSrvs []*backEnd.Server
    addrConfig string
}

func (cfg *Config) init (leng int){
    cfg.backEndSrvs = make([]*backEnd.Server, leng)
}

func (cfg *Config) initBackend(id int) {
    tmpsrv := backEnd.Server{}
    //cfg.mu.Lock()
    //defer cfg.mu.Unlock()
    StartBackEnd(id, &tmpsrv, cfg.addrConfig)
    cfg.backEndSrvs[id] =&tmpsrv
    //cfg.mu.Unlock()
}

func setUpAddress(srv *backEnd.Server, configFilePath string) {
    file, err := os.Open(configFilePath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        addressAndPort := scanner.Text()
        tokens := strings.Split(addressAndPort, ":")
        port := "80"
        if len(tokens) == 2 {
            port = tokens[1]
        }
        srv.RegisterAddress(tokens[0], port)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func StartBackEnd(id int, srv *backEnd.Server, config string) {
    fmt.Printf("StartBackEnd: %d \n", id)
    srv.Init(id)
    setUpAddress(srv, config)
    srv.Start()
    //return srv
}

func TestOneBackEndServer(t *testing.T){
    var frontEndSrv server.Server
    cfg := &Config{}
    addrConfig := "config.txt"
    frontEndSrv.InitialDial("tcp", addrConfig)
    addrBook := frontEndSrv.GetAddressBook()
    backEndNum := len(addrBook)
    cfg.init(backEndNum)
    cfg.addrConfig = addrConfig
    
    for id := 0; id < backEndNum; id++{
        fmt.Printf("id: %d\n", id)
        go cfg.initBackend(id)
        time.Sleep(2 * time.Second)
    }

    _, reply1 := handler.ClientRegisterUserRPC("User1", "user1pw", &frontEndSrv)
    token := reply1.Token
    handler.ClientPostRPC(token, "user1 post 1", &frontEndSrv)
    handler.ClientGetMyContentRPC(token, &frontEndSrv)
    handler.ClientPostRPC(token, "user1 post 2", &frontEndSrv)
    //for{

    //}
    if(true){
        t.Errorf("Token: %s.", token)
    }
}
