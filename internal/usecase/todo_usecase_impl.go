package usecase

import (
	"errors"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
)

type todoUsecase struct {
	repo domain.TodoRepository
}

// NewTodoUsecase is a Wire provider.
// Wire sees: needs domain.TodoRepository → returns TodoUsecase.
// It injects the repository that was already built in the infrastructure layer.
func NewTodoUsecase(repo domain.TodoRepository) TodoUsecase {
	return &todoUsecase{repo: repo}
}

func (u *todoUsecase) CreateTodo(title, description string) (*domain.Todo, error) {
	// Business rule: title is mandatory.
	// This validation belongs here, not in the handler (presentation) or repository (data).
	if title == "" {
		return nil, errors.New("title is required")
	}

	todo := &domain.Todo{
		Title:       title,
		Description: description,
	}

	return u.repo.Create(todo)
}

func (u *todoUsecase) GetTodo(id string) (*domain.Todo, error) {
	todo, err := u.repo.GetbyID(id)
	if err != nil {
		return nil, err
	}
	return todo, nil
}
