package backEnd

import (
    //"errors"
    "reflect"
    "backEnd/cmd"
    "fmt"
)

type Service struct {
    commands chan reflect.Value
}

func (service *Service) makeRPCHandler(args reflect.Value, reply reflect.Value, cmdType reflect.Type, message string) error {
    chanType := reflect.ChanOf(reflect.BothDir, reply.Type())
    channel := reflect.MakeChan(chanType, 1)
    command := reflect.Indirect(reflect.New(cmdType))
    command.Field(0).Set(args)
    command.Field(1).Set(channel)
    service.commands <- command
    result, _ := channel.Recv()
    fmt.Printf("%v\n", result)
    for i:=0; i<reply.Elem().NumField(); i++ {
        reply.Elem().Field(i).Set(result.Elem().Field(i))
    }
    fmt.Printf("%v\n", reply)
    fmt.Printf("%v", message)
    /*if ok && result.Elem().Field(0).Bool() {
        return nil
    }*/
    return nil
}

func (service *Service) ValidateUser(args cmd.ValidateUserArgs, reply *cmd.ValidateUserReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.ValidateUser{}),
        "ValidateUser error",
    )
}

func (service *Service) ValidateAuth(args cmd.ValidateAuthArgs, reply *cmd.ValidateAuthReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.ValidateAuth{}),
        "ValidateAuth error",
    )
}

func (service *Service) UserLogout(args cmd.UserLogoutArgs, reply *cmd.UserLogoutReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.UserLogout{}),
        "UserLogout error",
    )
}

func (service *Service) UserLogin(args cmd.UserLoginArgs, reply *cmd.UserLoginReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.UserLogin{}),
        "UserLogin error",
    )
}

func (service *Service) RegisterUser(args cmd.RegisterUserArgs, reply *cmd.RegisterUserReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(&args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.RegisterUser{}),
        "RegisterUser error",
    )
}

func (service *Service) DeleteUser(args cmd.DeleteUserArgs, reply *cmd.DeleteUserReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.DeleteUser{}),
        "DeleteUser error",
    )
}

func (service *Service) Post(args cmd.PostArgs, reply *cmd.PostReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.Post{}),
        "Post error",
    )
}

func (service *Service) Follow(args cmd.FollowArgs, reply *cmd.FollowReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.Follow{}),
        "Follow error",
    )
}

func (service *Service) UnFollow(args cmd.UnFollowArgs, reply *cmd.UnFollowReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.UnFollow{}),
        "UnFollow error",
    )
}

func (service *Service) GetMyContent(args cmd.GetMyContentArgs, reply *cmd.GetMyContentReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.GetMyContent{}),
        "GetMyContent error",
    )
}
