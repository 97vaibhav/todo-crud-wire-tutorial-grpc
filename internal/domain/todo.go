package domain

import "time"

// TodoStatus is a typed string so the compiler catches typos at usage sites.
type TodoStatus string

const (
	TodoStatusPending    TodoStatus = "PENDING"
	TodoStatusInProgress TodoStatus = "IN_PROGRESS"
	TodoStatusDone       TodoStatus = "DONE"
)

// Todo is the pure domain entity. No GORM tags, no proto tags — just Go.
// This struct travels through every layer of the application.
type Todo struct {
	ID          string
	Title       string
	Description string
	Status      TodoStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TodoRepository is the interface the usecase layer depends on.
// The infrastructure layer provides the concrete implementation.
// This is Dependency Inversion — high-level policy (usecase) does NOT depend
// on low-level detail (GORM). Both depend on this abstraction.
type TodoRepository interface {
	Create(todo *Todo) (*Todo, error)
}
