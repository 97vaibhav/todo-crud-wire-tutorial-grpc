package domain

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string // bcrypt hash — never expose outside this layer
	GroupID      string
	GroupType    GroupType // denormalised from the group for convenience
	CreatedAt    time.Time
}

type UserRepository interface {
	Create(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Delete(id string) error
	Update(id string, user *User) (*User, error)
}
