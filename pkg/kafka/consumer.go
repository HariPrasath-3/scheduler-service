package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type MessageHandler interface {
	HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type Consumer struct {
	group   sarama.ConsumerGroup
	handler MessageHandler
	topics  []string
}

func NewConsumer(
	brokers []string,
	groupID string,
	topics []string,
	handler MessageHandler,
) (*Consumer, error) {

	cfg := sarama.NewConfig()

	cfg.Version = sarama.V2_1_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange

	group, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		group:   group,
		handler: handler,
		topics:  topics,
	}, nil
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {
		if err := c.handler.HandleMessage(session.Context(), msg); err != nil {
			log.Printf("handler error: %v", err)
			// do NOT commit offset on error
			continue
		}

		session.MarkMessage(msg, "")
	}

	return nil
}

func (c *Consumer) Close() error {
	return c.group.Close()
}
