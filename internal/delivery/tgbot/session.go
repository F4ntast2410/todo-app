package tgbot

import "sync"

type UserState int

const (
	StateIdle UserState = iota
	StateWaitingTaskTitle
	StateWaitingTaskDescription
	StateWaitingNewTaskDescription
)

type UserSession struct {
	State     UserState
	TaskTitle string
	TaskID    int
}

type SessionCache struct {
	Cache map[int64]UserSession
	Mu    sync.RWMutex
}
