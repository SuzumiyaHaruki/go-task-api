package main

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
