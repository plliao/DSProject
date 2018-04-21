package backEnd

import (
    "backEnd/cmd"
)

type Service struct {
    commands chan Command
}

func (service *Service) ValidateUser(args cmd.ValidateUserArgs, reply *cmd.ValidateUserReply) error {
    channel := make(chan *cmd.ValidateUserReply, 1)
    cmd := cmd.ValidateUser{&args, channel}
    service.commands <- cmd
    reply = <-cmd.Channel
    return reply.Err
}
