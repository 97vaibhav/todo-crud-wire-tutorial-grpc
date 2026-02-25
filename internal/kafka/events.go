package kafka

import "time"

// Event type constants. Use these when publishing so consumers can filter or
// route by event_type. When you add Idea 2/3, add more event types here.
const (
	EventTypeTodoCreated = "todo.created"
	// EventTypeTodoUpdated = "todo.updated"  // for later
	// EventTypeTodoDeleted = "todo.deleted"  // for later
)

// TodoEvent is the payload we send to the todo-events topic.
// JSON-serialized so any consumer (audit, notification, analytics) can read it.
// Keep it flat and additive so new consumers don't break when you add fields.
type TodoEvent struct {
	EventType  string    `json:"event_type"`
	TodoID     string    `json:"todo_id"`
	UserID     string    `json:"user_id"`
	Title      string    `json:"title"`
	OccurredAt time.Time `json:"occurred_at"`
}
