package kafka

import (
	"context"

	"github.com/HariPrasath-3/scheduler-service/pkg/config"
	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(config *config.KafkaConfig) (*Producer, error) {
	cfg := sarama.NewConfig()

	// Required settings
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Return.Successes = true

	// Kafka version (safe default)
	cfg.Version = sarama.V2_1_0_0

	p, err := sarama.NewSyncProducer(config.Brokers, cfg)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Send(
	ctx context.Context,
	topic string,
	key string,
	value []byte,
) error {

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := p.producer.SendMessage(msg)
	return err
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
