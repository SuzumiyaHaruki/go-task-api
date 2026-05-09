package main

import (
	"net/http"
	"strings"
	"time"
)

func (a *app) listTasks(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	tasks := make([]task, 0, len(a.tasks))
	for _, item := range a.tasks {
		tasks = append(tasks, item)
	}

	writeOK(w, tasks)
}

func (a *app) createTask(w http.ResponseWriter, r *http.Request) {
	if !hasDemoToken(r) {
		writeError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
		return
	}

	var req taskRequest
	if !readJSON(w, r, &req) {
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	if req.Status == "" {
		req.Status = "todo"
	}

	now := time.Now()
	item := task{
		ID:        a.nextID(),
		Title:     req.Title,
		Content:   req.Content,
		Status:    req.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	a.mu.Lock()
	a.tasks[item.ID] = item
	a.mu.Unlock()

	writeCreated(w, item)
}

func (a *app) getTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	a.mu.Lock()
	item, exists := a.tasks[id]
	a.mu.Unlock()

	if !exists {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	writeOK(w, item)
}

func (a *app) updateTask(w http.ResponseWriter, r *http.Request) {
	if !hasDemoToken(r) {
		writeError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
		return
	}

	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req taskRequest
	if !readJSON(w, r, &req) {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	item, exists := a.tasks[id]
	if !exists {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	if strings.TrimSpace(req.Title) != "" {
		item.Title = strings.TrimSpace(req.Title)
	}
	if req.Content != "" {
		item.Content = req.Content
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	item.UpdatedAt = time.Now()
	a.tasks[id] = item

	writeOK(w, item)
}

func (a *app) deleteTask(w http.ResponseWriter, r *http.Request) {
	if !hasDemoToken(r) {
		writeError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
		return
	}

	id, ok := parseID(w, r)
	if !ok {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.tasks[id]; !exists {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	delete(a.tasks, id)

	writeOK(w, map[string]int64{"deleted_id": id})
}
