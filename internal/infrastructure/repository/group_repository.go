package repository

import (
	"time"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
	"gorm.io/gorm"
)

type groupModel struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	Type      string    `gorm:"not null"` // "ADMIN" or "GUEST"
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (groupModel) TableName() string { return "groups" }

func (m *groupModel) toDomain() *domain.Group {
	return &domain.Group{
		ID:        m.ID,
		Name:      m.Name,
		Type:      domain.GroupType(m.Type),
		CreatedAt: m.CreatedAt,
	}
}

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) domain.GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) FindByID(id string) (*domain.Group, error) {
	var model groupModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.toDomain(), nil
}

func (r *groupRepository) List() ([]*domain.Group, error) {
	var models []groupModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	groups := make([]*domain.Group, len(models))
	for i, m := range models {
		g := m
		groups[i] = g.toDomain()
	}
	return groups, nil
}
