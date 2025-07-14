package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/JakeTheDoggg/taskmanager/internal/handlers"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	// http.HandleFunc("/tasks/create", handlers.CreateTaskHandler)
	// http.HandleFunc("/tasks/get", handlers.GetTaskHandler)
	// http.HandleFunc("/tasks/getall", handlers.GetAllTasksHandler)

	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetAllTasksHandler(w, r)
		case http.MethodPost:
			handlers.CreateTaskHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodPut:
			handlers.UpdateTaskHandler(w, r, id)
		case http.MethodDelete:
			handlers.DeleteTaskHandler(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
