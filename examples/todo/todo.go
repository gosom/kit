package todo

import "fmt"

type Todo struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func NewTodo(id string) Todo {
	return Todo{ID: id}
}

func (t *Todo) String() string {
	return fmt.Sprintf("Todo { id: %s, title: %s, status: %s }", t.ID, t.Title, t.Status)
}

func (o *Todo) UpdateStatus(s string) error {
	switch s {
	case "completed":
		if o.Status == "open" {
			o.Status = "completed"
		} else {
			return fmt.Errorf("invalid status transition from %s to %s", o.Status, s)
		}
	case "open":
		if o.Status == "completed" {
			o.Status = "open"
		} else {
			return fmt.Errorf("invalid status transition from %s to %s", o.Status, s)
		}
	default:
		return fmt.Errorf("invalid status %s", s)
	}
	return nil
}
