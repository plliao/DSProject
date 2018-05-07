package backEnd

import (
    "backEnd/cmd"
)

func (srv *Server) UserLogout(token string) (bool, string) {
    ok := srv.validateAuth(token)
    if ok {
        user := srv.tokens[token]
        srv.deleteUserToken(user)
        return true, srv.messages.NoError
    }
    return false, srv.messages.UnrecognizedToken
}

func (srv *Server) UserLogin(username string, password string) (bool, string, string) {
    ok, err := srv.validateUser(username, password)
    if ok {
        user := srv.users[username]
        srv.generateUserToken(user)
        return true, srv.messages.NoError, user.token
    }
    return false, err.Error(), srv.messages.EmptyToken
}

func (srv *Server) RegisterUser(username string, password string) (bool, string, string) {
    if ok, err := srv.validateUserNameAndPassFormat(username, password); !ok {
        return ok, err.Error(), srv.messages.EmptyToken
    }
    if _, ok := srv.users[username]; ok {
        return false, srv.messages.UserAlreadyExist, srv.messages.EmptyToken
    }
    newUser := &User{
        Username:username,
        Password:password,
    }
    newUser.Init()
    srv.users[username] = newUser
    srv.generateUserToken(newUser)
    return true, srv.messages.NoError, newUser.token
}

func (srv *Server) DeleteUser(token string) (bool, string) {
    ok := srv.validateAuth(token)
    if ok {
        user := srv.tokens[token]
        srv.removeUser(user)
        return true, srv.messages.NoError
    }
    return false, srv.messages.UnrecognizedToken
}

func (srv *Server) Post(token string, content string) (bool, string) {
    ok := srv.validateAuth(token)
    if ok {
        user := srv.tokens[token]
        user.Post(content)
        return true, srv.messages.NoError
    }
    return false, srv.messages.UnrecognizedToken
}

func (srv *Server) Follow(token string, username string) (bool, string) {
    ok := srv.validateAuth(token)
    if ok {
        follower := srv.tokens[token]
        if user, hasUser := srv.users[username]; hasUser {
            follower.Follow(user)
            return true, srv.messages.NoError
        }
        return false, username + " " + srv.messages.UserNotExist
    }
    return false, srv.messages.UnrecognizedToken
}

func (srv *Server) UnFollow(token string, username string) (bool, string) {
    ok := srv.validateAuth(token)
    if ok {
        follower := srv.tokens[token]
        if user, hasUser := srv.users[username]; hasUser {
            follower.UnFollow(user)
            return true, srv.messages.NoError
        }
        return false, username + " " + srv.messages.UserNotExist
    }
    return false, srv.messages.UnrecognizedToken
}

func (srv *Server) GetMyContent(token string) (bool, string, []*cmd.Article) {
    ok := srv.validateAuth(token)
    if ok {
        user := srv.tokens[token]
        return true, srv.messages.NoError, user.GetMyContent()
    }
    return false, srv.messages.UnrecognizedToken, nil
}

func (srv *Server) GetFollower(token string) (bool, string, []*cmd.Relationship) {
    ok := srv.validateAuth(token)
    if ok {
        relationships := make([]*cmd.Relationship, 0, len(srv.users))
        follower := srv.tokens[token]
        for username, _ := range srv.users {
            if username != follower.Username {
                _, isFollowing := follower.following[username]
                relationships = append(relationships, &cmd.Relationship{username, isFollowing})
            }
        }
        return true, srv.messages.NoError, relationships
    }
    return false, srv.messages.UnrecognizedToken, nil
}
