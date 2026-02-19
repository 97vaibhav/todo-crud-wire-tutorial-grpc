package domain

import "time"

// GroupType controls what a user is permitted to do.
// ADMIN  → full user management rights
// GUEST  → can only manage their own todos
type GroupType string

const (
	GroupTypeAdmin GroupType = "ADMIN"
	GroupTypeGuest GroupType = "GUEST"
)

type Group struct {
	ID        string
	Name      string
	Type      GroupType
	CreatedAt time.Time
}

type GroupRepository interface {
	FindByID(id string) (*Group, error)
	List() ([]*Group, error)
}
