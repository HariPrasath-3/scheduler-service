#!/usr/bin/env sh

set -e

KAFKA_BROKER="kafka:9092"

echo "[kafka-init] Waiting for Kafka..."

until kafka-topics --bootstrap-server "$KAFKA_BROKER" --list >/dev/null 2>&1; do
  sleep 2
done

echo "[kafka-init] Kafka is ready"

create_topic() {
  TOPIC_NAME=$1
  PARTITIONS=$2
  REPLICATION=$3

  echo "[kafka-init] Ensuring topic: $TOPIC_NAME"

  kafka-topics \
    --bootstrap-server "$KAFKA_BROKER" \
    --create \
    --if-not-exists \
    --topic "$TOPIC_NAME" \
    --partitions "$PARTITIONS" \
    --replication-factor "$REPLICATION"
}

# Scheduler topics
create_topic "schedule_events" 4 1

echo "[kafka-init] Topic initialization complete"