package repository

import (
	"time"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// auditLogModel is the GORM representation of an audit_logs row.
type auditLogModel struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	EventType string    `gorm:"column:event_type;not null"`
	TodoID    string    `gorm:"column:todo_id;type:uuid"`
	UserID    string    `gorm:"column:user_id;type:uuid"`
	Payload   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (auditLogModel) TableName() string {
	return "audit_logs"
}

type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository is a Wire provider. Returns domain.AuditLogRepository.
func NewAuditLogRepository(db *gorm.DB) domain.AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *domain.AuditLog) error {
	model := &auditLogModel{
		ID:        uuid.New().String(),
		EventType: log.EventType,
		TodoID:    log.TodoID,
		UserID:    log.UserID,
		Payload:   log.Payload,
		CreatedAt: log.CreatedAt,
	}
	if model.CreatedAt.IsZero() {
		model.CreatedAt = time.Now().UTC()
	}
	return r.db.Create(model).Error
}
