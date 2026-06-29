package tgbot

import "sync"

type UserState int

const (
	StateIdle UserState = iota
	StateWaitingTaskTitle
	StateWaitingTaskDescription
)

type UserSession struct {
	State     UserState
	TaskTitle string
}

type SessionCache struct {
	Cache map[int64]UserSession
	Mu    sync.RWMutex
}
