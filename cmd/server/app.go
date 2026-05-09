package main

import (
	"net/http"
	"sync"
)

type app struct {
	mux         *http.ServeMux
	mu          sync.Mutex
	nextUserID  int64
	nextTaskID  int64
	usersByName map[string]user
	tasks       map[int64]task
}

func newApp() *app {
	a := &app{
		mux:         http.NewServeMux(),
		nextUserID:  1,
		nextTaskID:  1,
		usersByName: make(map[string]user),
		tasks:       make(map[int64]task),
	}

	a.routes()
	return a
}

func (a *app) nextID() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	id := a.nextTaskID
	a.nextTaskID++
	return id
}
