package repository

import (
	"time"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// todoModel is the GORM representation of a todo row in Postgres.
// It lives ONLY in this file — no other layer sees it.
// Notice: field types match the DB column types exactly.
type todoModel struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Status      string    `gorm:"not null;default:'PENDING'"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TableName tells GORM which table to use. Without this, GORM would guess "todo_models".
func (todoModel) TableName() string {
	return "todos"
}

// toDomain converts the infrastructure model → domain entity.
// This is the boundary crossing point.
func (m *todoModel) toDomain() *domain.Todo {
	return &domain.Todo{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Status:      domain.TodoStatus(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// todoRepository is the concrete implementation of domain.TodoRepository.
// It holds a *gorm.DB — the only place in the app that does so.
type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository is a Wire provider.
// It returns the interface type (domain.TodoRepository), NOT the concrete type.
// This enforces that the rest of the app only knows about the interface.
func NewTodoRepository(db *gorm.DB) domain.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *domain.Todo) (*domain.Todo, error) {
	model := &todoModel{
		ID:          uuid.New().String(),
		Title:       todo.Title,
		Description: todo.Description,
		Status:      string(domain.TodoStatusPending),
	}

	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}

	return model.toDomain(), nil
}
