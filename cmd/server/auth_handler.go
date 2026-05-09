package main

import (
	"net/http"
	"strconv"
	"strings"
)

func (a *app) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if !readJSON(w, r, &req) {
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "username is required and password must be at least 6 characters")
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.usersByName[req.Username]; exists {
		writeError(w, http.StatusConflict, "username already exists")
		return
	}

	u := user{
		ID:       a.nextUserID,
		Username: req.Username,
		Password: req.Password,
	}
	a.nextUserID++
	a.usersByName[u.Username] = u

	writeCreated(w, u)
}

func (a *app) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !readJSON(w, r, &req) {
		return
	}

	a.mu.Lock()
	u, exists := a.usersByName[req.Username]
	a.mu.Unlock()

	if !exists || u.Password != req.Password {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	writeOK(w, map[string]string{
		"token": "demo-token-" + strconv.FormatInt(u.ID, 10),
		"type":  "Bearer",
	})
}
