package common

const (
	KafkaTopicScheduleEvents = "schedule_events"
)

const (
	// RedisKeyFormatterScheduledEventIDs is used to store scheduled event IDs
	// in a specific bucket and partition.
	RedisKeyFormatterScheduledEvents = "scheduler:scheduled:bucket:%d:partition:%d"

	// RedisKeyFormatterProcessingEventIDs is used to store processing event IDs
	RedisKeyFormatterProcessingEvents = "scheduler:processing:bucket:%d:partition:%d"
)
