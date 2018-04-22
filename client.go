package main

import (
	"fmt"
	"log"
	"net/rpc"
    "backEnd/cmd"
    "strconv"
    "time"
    "math/rand"
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

    args.Username = "NotFollow"
    method = "RegisterUser"
    err = client.Call("Service." + method, args, &reply)
	fmt.Printf(method + ": +%v\n", reply)

    args.Username = "Following"
    method = "RegisterUser"
    err = client.Call("Service." + method, args, &reply)
	fmt.Printf(method + ": +%v\n", reply)

	loginArgs := cmd.UserLoginArgs{"plliao1234", "abcdefrg"}
	loginReply := cmd.UserLoginReply{}
    method = "UserLogin"
    err = client.Call("Service." + method, loginArgs, &loginReply)
	if err != nil {
		log.Fatal(method + " Error:", err)
	}
	fmt.Printf(method + ": +%v\n", loginReply)
    token := loginReply.Token

    followArgs := cmd.FollowArgs{token, args.Username}
    followReply := cmd.FollowReply{}
    method = "Follow"
    err = client.Call("Service." + method, followArgs, &followReply)
	if err != nil {
		log.Fatal(method + " Error:", err)
	}
	fmt.Printf(method + ": +%v\n", followReply)

    getFollowerArgs := cmd.GetFollowerArgs{token}
    getFollowerReply := cmd.GetFollowerReply{}
    method = "GetFollower"
    err = client.Call("Service." + method, getFollowerArgs, &getFollowerReply)
	if err != nil {
		log.Fatal(method + " Error:", err)
	}
	fmt.Printf(method + ": +%v\n", getFollowerReply)
    for i,p := range getFollowerReply.Relationships {
        fmt.Printf("User %v: %v\n", i, *p)
    }

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

    username := "user"
    password := "password"
    for i:=0; i<1000; i++ {
        go clientRead(username + strconv.Itoa(i), password)
    }
    time.Sleep(1000 * time.Second)
}

func clientRead(username string, password string) {
    client, err := rpc.DialHTTP("tcp", ":8080")
    if err != nil {
        log.Fatal("dialing:", err)
    }

	args := cmd.RegisterUserArgs{username, password}
	reply := cmd.RegisterUserReply{}
    method := "RegisterUser"
    err = client.Call("Service." + method, args, &reply)
	fmt.Printf(method + ": +%v\n", reply)
    token := reply.Token

    for i:=0; i>=0; i++ {
        myContentArgs := cmd.GetMyContentArgs{token}
        myContentReply := cmd.GetMyContentReply{}
        method = "GetMyContent"
        err = client.Call("Service." + method, myContentArgs, &myContentReply)
        //fmt.Printf(method + ": +%v\n", myContentReply)
        for i,p := range myContentReply.Articles {
            fmt.Printf(username + " Article %v: %v\n", i, *p)
        }

        if rand.Int31n(100) > 20 {
            articleId := len(myContentReply.Articles)
            content := username + " " + strconv.Itoa(articleId) + " post"
            postArgs := cmd.PostArgs{token, content}
            postReply := cmd.PostReply{}
            method = "Post"
            err = client.Call("Service." + method, postArgs, &postReply)
	        //fmt.Printf(method + ": +%v\n", postReply)
            fmt.Printf(username + " post " + content + "\n")
        }
        time.Sleep(time.Duration(rand.Int31n(300)) * time.Millisecond)

    }
}



