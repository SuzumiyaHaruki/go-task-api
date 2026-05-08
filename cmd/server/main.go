package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type apiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type taskRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type user struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type app struct {
	mux          *http.ServeMux
	mu           sync.Mutex
	nextUserID   int64
	nextTaskID   int64
	usersByName  map[string]user
	tasks        map[int64]task
	demoAuthUser string
}

func main() {
	port := getenv("APP_PORT", "8080")

	a := newApp()
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      logRequests(a.mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("go-task-api listening on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func newApp() *app {
	a := &app{
		mux:          http.NewServeMux(),
		nextUserID:   1,
		nextTaskID:   1,
		usersByName:  make(map[string]user),
		tasks:        make(map[int64]task),
		demoAuthUser: "demo-user",
	}

	a.routes()
	return a
}

func (a *app) routes() {
	a.mux.HandleFunc("GET /health", a.health)
	a.mux.HandleFunc("POST /api/v1/auth/register", a.register)
	a.mux.HandleFunc("POST /api/v1/auth/login", a.login)
	a.mux.HandleFunc("GET /api/v1/tasks", a.listTasks)
	a.mux.HandleFunc("POST /api/v1/tasks", a.createTask)
	a.mux.HandleFunc("GET /api/v1/tasks/{id}", a.getTask)
	a.mux.HandleFunc("PUT /api/v1/tasks/{id}", a.updateTask)
	a.mux.HandleFunc("DELETE /api/v1/tasks/{id}", a.deleteTask)
}

func (a *app) health(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]string{"status": "ok"})
}

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

func (a *app) nextID() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	id := a.nextTaskID
	a.nextTaskID++
	return id
}

func readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return false
	}

	return true
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}

	return id, true
}

func hasDemoToken(r *http.Request) bool {
	value := r.Header.Get("Authorization")
	return strings.HasPrefix(value, "Bearer demo-token-")
}

func writeOK(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, apiResponse{Code: 0, Message: "ok", Data: data})
}

func writeCreated(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusCreated, apiResponse{Code: 0, Message: "created", Data: data})
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, apiResponse{Code: status, Message: message})
}

func writeJSON(w http.ResponseWriter, status int, resp apiResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("write response failed: %v", err)
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
