package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/kafka"
)

type todoUsecase struct {
	repo     domain.TodoRepository
	producer kafka.Producer
}

// NewTodoUsecase is a Wire provider.
func NewTodoUsecase(repo domain.TodoRepository, producer kafka.Producer) TodoUsecase {
	return &todoUsecase{repo: repo, producer: producer}
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

	created, err := u.repo.Create(todo)
	if err != nil {
		return nil, err
	}

	// Publish event for audit (and later: notifications, analytics).
	// We use context.Background() so the publish is not cancelled when the gRPC request ends.
	ev := kafka.TodoEvent{
		EventType:  kafka.EventTypeTodoCreated,
		TodoID:     created.ID,
		UserID:     created.UserID,
		Title:      created.Title,
		OccurredAt: created.CreatedAt,
	}
	payload, _ := json.Marshal(ev)
	if err := u.producer.Publish(context.Background(), kafka.TopicTodoEvents, payload); err != nil {
		log.Printf("[todo-usecase] kafka publish: %v", err)
		// Do not fail the request; audit is best-effort.
	}

	return created, nil
}

func (u *todoUsecase) ListTodos() ([]*domain.Todo, error) {
	return u.repo.List()
}

func (u *todoUsecase) GetTodo(id string) (*domain.Todo, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return u.repo.GetbyID(id)
}
