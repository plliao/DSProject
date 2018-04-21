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
    method := "RegisterUser"
    err = client.Call("Service." + method, args, &reply)
	fmt.Printf(method + ": +%v\n", reply)

	loginArgs := cmd.UserLoginArgs{"plliao1234", "abcdefrg"}
	loginReply := cmd.UserLoginReply{}
    method = "UserLogin"
    err = client.Call("Service." + method, loginArgs, &loginReply)
	if err != nil {
		log.Fatal(method + " Error:", err)
	}
	fmt.Printf(method + ": +%v\n", &loginReply)


    token := loginReply.Token
    content := "My first post"
    postArgs := cmd.PostArgs{token, content}
    postReply := cmd.PostReply{}
    method = "Post"
    err = client.Call("Service." + method, postArgs, &postReply)
	if err != nil {
		log.Fatal(method + " Error:", err)
	}
	fmt.Printf(method + ": +%v\n", postReply)

    myContentArgs := cmd.GetMyContentArgs{token}
    myContentReply := cmd.GetMyContentReply{}
    method = "GetMyContent"
    err = client.Call("Service." + method, myContentArgs, &myContentReply)
	if err != nil {
		log.Fatal(method + " Error:", err)
	}
    fmt.Printf(method + ": +%v\n", myContentReply)
    for i,p := range myContentReply.Articles {
        fmt.Printf("Article %v: %v\n", i, *p)
    }
}



