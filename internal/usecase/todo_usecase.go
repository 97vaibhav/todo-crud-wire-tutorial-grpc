package usecase

import "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"

// TodoUsecase defines what operations the application supports.
// The gRPC handler depends on this interface, NOT the concrete implementation.
// When you write tests for the handler, you mock this — not GORM.
type TodoUsecase interface {
	CreateTodo(title, description string) (*domain.Todo, error)
	GetTodo(id string) (*domain.Todo, error)
}
