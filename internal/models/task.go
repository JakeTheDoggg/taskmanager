package models

import "fmt"

type Task struct {
	ID          int    `json:"id"`          // Обратите внимание на теги!
	Title       string `json:"title"`       // Должно совпадать с
	Description string `json:"description"` // названиями полей в JSON
	IsCompleted bool   `json:"isCompleted"` // camelCase важно!
}

func (t Task) String() string {
	var status string
	if t.IsCompleted {
		status = "DONE"
	} else {
		status = "TODO"
	}
	return fmt.Sprintf("Task #%v.\nTitle:%v.\n%v\nStatus: %v", t.ID, t.Title, t.Description, status)
}
