package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/JakeTheDoggg/taskmanager/internal/models"
	"github.com/JakeTheDoggg/taskmanager/internal/storage"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса ДВАЖДЫ (для логов и обработки)
	body, _ := io.ReadAll(r.Body)
	fmt.Printf("Raw body: %s\n", string(body))

	// Важно: создаем новый Reader после чтения
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		fmt.Printf("Decode error: %v\n", err) // Логируем ошибку
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed task: %+v\n", task) // Лог распарсенных данных

	// Добавьте эту проверку
	if task.Title == "" {
		errMsg := "Title is required"
		fmt.Println(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	id, err := storage.AddTask(task)
	if err != nil {
		fmt.Printf("AddTask error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := storage.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks := storage.GetAllTasks()
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("JSON encode error: %v", err)
	}
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	var updateData struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		IsCompleted *bool   `json:"isCompleted"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, err := storage.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if updateData.Title != nil {
		task.Title = *updateData.Title
	}
	if updateData.Description != nil {
		task.Description = *updateData.Description
	}
	if updateData.IsCompleted != nil {
		task.IsCompleted = *updateData.IsCompleted
	}

	if err := storage.UpdateTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := storage.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
func DeleteAllTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	storage.DeleteAllTasks()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "all tasks deleted"})
}
