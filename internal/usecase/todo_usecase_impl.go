package usecase

import (
	"errors"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
)

type todoUsecase struct {
	repo domain.TodoRepository
}

// NewTodoUsecase is a Wire provider.
func NewTodoUsecase(repo domain.TodoRepository) TodoUsecase {
	return &todoUsecase{repo: repo}
}

func (u *todoUsecase) CreateTodo(userID, title, description string) (*domain.Todo, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	todo := &domain.Todo{
		Title:       title,
		Description: description,
		UserID:      userID, // comes from the JWT, not the request body
	}

	return u.repo.Create(todo)
}
