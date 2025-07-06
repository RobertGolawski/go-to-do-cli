package models

import (
	"time"
)

type Priority int

const (
	Low Priority = iota
	Medium
	High
	Critical
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

type TodoItem struct {
	ID       int       `json:"id"`
	ToDo     string    `json:"todo"`
	DueDate  time.Time `json:"due_date"`
	Priority Priority  `json:"priority"`
	Done     bool      `json:"done"`
}

type TodoList struct {
	Todos  []TodoItem `json:"items"`
	NextID int        `json:"next_id"`
}

func (tl *TodoList) AddTodo(todo string, dueDate time.Time, priority Priority) {
	t := TodoItem{
		ID:       tl.NextID,
		ToDo:     todo,
		DueDate:  dueDate,
		Priority: priority,
		Done:     false,
	}
	tl.Todos = append(tl.Todos, t)
	tl.NextID++
}
