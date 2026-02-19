package usecase

import "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"

type AuthUsecase interface {
	// Login verifies credentials and returns a signed JWT on success.
	Login(email, password string) (token string, user *domain.User, err error)

	// CreateUser is called by admins to provision a new user in a given group.
	CreateUser(name, email, password, groupID string) (*domain.User, error)

	// DeleteUser is called by admins to remove a user by ID.
	DeleteUser(userID string) error

	// ListGroups returns all groups so admins can pick the right group_id.
	ListGroups() ([]*domain.Group, error)

	UpdateUser(id string, user *domain.User) (*domain.User, error)
}
