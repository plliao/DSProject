package main

import (
    "flag"
    "backEnd"
    "strconv"
)

func main() {
    port := flag.Int("port", 8080, "Serving port")
    flag.Parse()

    var srv backEnd.Server
    srv.Init()
    srv.Start(strconv.Itoa(*port))
}
