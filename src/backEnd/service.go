package backEnd

import (
    "reflect"
    "backEnd/cmd"
)

type Service struct {
    commands chan reflect.Value
}

func (service *Service) makeRPCHandler(args reflect.Value, reply reflect.Value, cmdType reflect.Type) error {
    chanType := reflect.ChanOf(reflect.BothDir, reply.Type())
    channel := reflect.MakeChan(chanType, 1)
    command := reflect.Indirect(reflect.New(cmdType))
    command.Field(0).Set(args)
    command.Field(1).Set(channel)
    service.commands <- command
    result, _ := channel.Recv()
    for i:=0; i<reply.Elem().NumField(); i++ {
        reply.Elem().Field(i).Set(result.Elem().Field(i))
    }
    return nil
}

func (service *Service) UserLogout(args cmd.UserLogoutArgs, reply *cmd.UserLogoutReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(&args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.UserLogout{}),
    )
}

func (service *Service) UserLogin(args cmd.UserLoginArgs, reply *cmd.UserLoginReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(&args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.UserLogin{}),
    )
}

func (service *Service) RegisterUser(args cmd.RegisterUserArgs, reply *cmd.RegisterUserReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(&args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.RegisterUser{}),
    )
}

func (service *Service) DeleteUser(args cmd.DeleteUserArgs, reply *cmd.DeleteUserReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(&args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.DeleteUser{}),
    )
}

func (service *Service) Post(args cmd.PostArgs, reply *cmd.PostReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(&args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.Post{}),
    )
}

func (service *Service) Follow(args cmd.FollowArgs, reply *cmd.FollowReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(&args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.Follow{}),
    )
}

func (service *Service) UnFollow(args cmd.UnFollowArgs, reply *cmd.UnFollowReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(&args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.UnFollow{}),
    )
}

func (service *Service) GetMyContent(args cmd.GetMyContentArgs, reply *cmd.GetMyContentReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(&args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.GetMyContent{}),
    )
}

func (service *Service) GetFollower(args cmd.GetFollowerArgs, reply *cmd.GetFollowerReply) error {
    return service.makeRPCHandler(
        reflect.ValueOf(&args),
        reflect.ValueOf(reply),
        reflect.TypeOf(cmd.GetFollower{}),
    )
}

