package repository

import (
	"time"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userModel struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	GroupID      string    `gorm:"type:uuid;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (userModel) TableName() string { return "users" }

func (m *userModel) toDomain(groupType domain.GroupType) *domain.User {
	return &domain.User{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		GroupID:      m.GroupID,
		GroupType:    groupType,
		CreatedAt:    m.CreatedAt,
	}
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) (*domain.User, error) {
	model := &userModel{
		ID:           uuid.New().String(),
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		GroupID:      user.GroupID,
	}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}
	return model.toDomain(user.GroupType), nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var model userModel
	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}
	// Join the group to get its type.
	var gm groupModel
	if err := r.db.Where("id = ?", model.GroupID).First(&gm).Error; err != nil {
		return nil, err
	}
	return model.toDomain(domain.GroupType(gm.Type)), nil
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var model userModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	var gm groupModel
	if err := r.db.Where("id = ?", model.GroupID).First(&gm).Error; err != nil {
		return nil, err
	}
	return model.toDomain(domain.GroupType(gm.Type)), nil
}

func (r *userRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&userModel{}).Error
}
