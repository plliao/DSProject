package server

import (
	user "../user"
)

type Server struct {
    Users map[string]*user.User
}
