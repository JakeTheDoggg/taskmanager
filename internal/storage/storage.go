package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/JakeTheDoggg/taskmanager/internal/models"
)

var (
	tasks  = make(map[int]models.Task)
	lastID int
	mutex  sync.Mutex
)

func AddTask(t models.Task) (int, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if t.Title == "" {
		return 0, errors.New("title cannot be empty")
	}
	lastID++
	t.ID = lastID
	tasks[lastID] = t
	return lastID, nil

}

func GetTask(id int) (models.Task, error) {
	task, ok := tasks[id]
	if !ok {
		return models.Task{}, fmt.Errorf("task with ID %d not found", id)
	}
	return task, nil
}

func GetAllTasks() []models.Task {
	allTasks := make([]models.Task, 0, len(tasks))
	for _, task := range tasks {
		allTasks = append(allTasks, task)
	}
	return allTasks
}

func UpdateTask(t models.Task) error {
	if _, exists := tasks[t.ID]; !exists {
		return fmt.Errorf("task with ID %d not found", t.ID)
	}
	if t.Title == "" {
		return errors.New("title cannot be empty")
	}
	tasks[t.ID] = t
	return nil
}

func DeleteTask(id int) error {
	if _, exists := tasks[id]; !exists {
		return fmt.Errorf("task with ID %d not found", id)
	}
	delete(tasks, id)
	return nil
}
func DeleteAllTasks() {
	mutex.Lock()
	defer mutex.Unlock()
	tasks = make(map[int]models.Task)
	lastID = 0
}
