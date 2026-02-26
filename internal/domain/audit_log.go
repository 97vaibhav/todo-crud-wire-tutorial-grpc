package domain

import "time"

// AuditLog is the domain entity for an audit log entry.
// The consumer writes these after reading events from Kafka.
type AuditLog struct {
	ID        string
	EventType string
	TodoID    string
	UserID    string
	Payload   string // raw JSON for flexibility; consumers can store the full event
	CreatedAt time.Time
}

// AuditLogRepository is the interface the audit consumer depends on.
// The infrastructure layer provides the GORM implementation.
type AuditLogRepository interface {
	Create(log *AuditLog) error
}
