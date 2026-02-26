package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/config"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
)

// Producer is the interface the usecase layer depends on.
// The usecase only calls Publish; it does not depend on franz-go or any
// concrete Kafka client. Later, Idea 2/3 can reuse the same producer.
type Producer interface {
	// Publish sends a message to the given topic. Payload is the raw bytes
	// (e.g. JSON). Non-blocking best-effort: we log errors but do not fail
	// the request if Kafka is down, so the app stays resilient.
	Publish(ctx context.Context, topic string, payload []byte) error
}

// producerImpl wraps franz-go client and implements Producer.
type producerImpl struct {
	client *kgo.Client
}

// NewProducer is a Wire provider. It builds a Kafka client, ensures the
// todo-events topic exists, and returns the Producer interface.
func NewProducer(cfg *config.Config) (Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.KafkaBroker),
	)
	if err != nil {
		return nil, err
	}
	// Create the topic if it doesn't exist (Redpanda/Kafka may not auto-create).
	if err := ensureTopic(context.Background(), client, TopicTodoEvents); err != nil {
		log.Printf("[kafka] ensure topic %q: %v", TopicTodoEvents, err)
		// Continue anyway; produce will fail with a clear error if topic is missing.
	}
	return &producerImpl{client: client}, nil
}

// ensureTopic creates the topic if it does not exist. Idempotent: safe to call at startup.
func ensureTopic(ctx context.Context, client *kgo.Client, topic string) error {
	req := kmsg.NewPtrCreateTopicsRequest()
	req.Topics = []kmsg.CreateTopicsRequestTopic{{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}}
	resp, err := client.Request(ctx, req)
	if err != nil {
		return err
	}
	createResp := resp.(*kmsg.CreateTopicsResponse)
	for _, t := range createResp.Topics {
		if t.Topic != topic {
			continue
		}
		// 0 = success, 36 = TOPIC_ALREADY_EXISTS (idempotent)
		if t.ErrorCode != 0 && t.ErrorCode != 36 {
			msg := ""
			if t.ErrorMessage != nil {
				msg = *t.ErrorMessage
			}
			return fmt.Errorf("create topic %q: code %d %s", topic, t.ErrorCode, msg)
		}
		break
	}
	return nil
}

func (p *producerImpl) Publish(ctx context.Context, topic string, payload []byte) error {
	result := p.client.ProduceSync(ctx, &kgo.Record{
		Topic: topic,
		Value: payload,
	})
	if result.FirstErr() != nil {
		return result.FirstErr()
	}
	return nil
}
