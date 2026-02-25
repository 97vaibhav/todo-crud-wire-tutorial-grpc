package kafka

// Topic names live in one place so producers and consumers stay in sync.
// When you add Idea 2 (notifications) or Idea 3 (analytics), you can either
// reuse this topic (multiple consumer groups) or add new topic constants here.
const (
	// TopicTodoEvents is the topic for todo-related events (create, update, etc.).
	// The audit consumer reads from this topic. Later, a notification consumer
	// can read from the same topic with a different consumer group.
	TopicTodoEvents = "todo-events"
)
