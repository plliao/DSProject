package backend

import (
    "reflect"
)

type Raft struct {
    term int
    commitIndex int
    index int

    logs []reflect.Value
}
