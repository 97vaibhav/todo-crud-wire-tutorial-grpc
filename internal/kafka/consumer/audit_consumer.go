package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/config"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	// ConsumerGroupAudit is the consumer group for the audit log consumer.
	// When you add Idea 2 (notifications), use a different group (e.g. "todo-notifications")
	// so both consumers get the same events independently.
	ConsumerGroupAudit = "todo-audit"
)

// AuditConsumer reads from the todo-events topic and writes audit log entries to the DB.
// Run blocks until ctx is cancelled; start it in a goroutine from main.
type AuditConsumer struct {
	client *kgo.Client
	repo   domain.AuditLogRepository
}

// NewAuditConsumer is a Wire provider. It builds a Kafka consumer client
// subscribed to todo-events with consumer group "todo-audit".
func NewAuditConsumer(cfg *config.Config, repo domain.AuditLogRepository) (*AuditConsumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.KafkaBroker),
		kgo.ConsumerGroup(ConsumerGroupAudit),
		kgo.ConsumeTopics(kafka.TopicTodoEvents),
	)
	if err != nil {
		return nil, err
	}
	return &AuditConsumer{client: client, repo: repo}, nil
}

// Run processes records from the todo-events topic and persists them as audit logs.
// Call this in a goroutine; it runs until ctx is cancelled.
func (c *AuditConsumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.client.Close()
			return
		default:
			fetches := c.client.PollFetches(ctx)
			if fetches.IsClientClosed() {
				return
			}
			fetches.EachRecord(func(rec *kgo.Record) {
				if err := c.processRecord(ctx, rec); err != nil {
					log.Printf("[audit-consumer] process record: %v", err)
					// We still commit so we don't reprocess forever on bad data.
					// In production you might send to a dead-letter topic.
				}
			})
		}
	}
}

func (c *AuditConsumer) processRecord(ctx context.Context, rec *kgo.Record) error {
	var ev kafka.TodoEvent
	if err := json.Unmarshal(rec.Value, &ev); err != nil {
		return err
	}
	payload := string(rec.Value)
	auditEntry := &domain.AuditLog{
		EventType: ev.EventType,
		TodoID:    ev.TodoID,
		UserID:    ev.UserID,
		Payload:   payload,
		CreatedAt: time.Now().UTC(),
	}
	return c.repo.Create(auditEntry)
}
