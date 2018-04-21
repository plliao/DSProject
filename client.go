package main

import (
	"fmt"
	"log"
	"net/rpc"
    "backEnd/cmd"

)

func main() {
    client, err := rpc.DialHTTP("tcp", ":8080")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args := cmd.RegisterUserArgs{"plliao1234", "abcdefrg"}
	reply := cmd.RegisterUserReply{}
	err = client.Call("Service.RegisterUser", args, &reply)
	if err != nil {
		log.Fatal("registeration error:", err)
	}
	fmt.Printf("registeration: +%v\n", reply)
}
