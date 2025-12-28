package common

const (
	KafkaTopicScheduleEvents = "schedule_events"
)

const (
	// RedisKeyFormatterScheduledEventIDs is used to store scheduled event IDs
	// in a specific bucket and partition.
	RedisKeyFormatterScheduledEvents = "scheduled:bucket:%d:partition:%d"

	RedisKeyFormatterProcessingEvents = "processing:bucket:%d:partition:%d"
)
