package main

import (
    "flag"
    "backEnd"
    "strings"
    "bufio"
    "os"
    "log"
)

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

func main() {
    id := flag.Int("id", 0, "Server id in config")
    config := flag.String("config", "config.txt", "Replica public address config file")
    flag.Parse()

    var srv backEnd.Server
    srv.Init(*id)
    setUpAddress(&srv, *config)
    srv.Start()
}
