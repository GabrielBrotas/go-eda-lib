package kafka

import (
	"log"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer *ckafka.Consumer
}

func NewConsumer(configMap *ckafka.ConfigMap, topics []string) (*Consumer, error) {
	c, err := ckafka.NewConsumer(configMap)
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: c}, nil
}

// Consume consumes messages from Kafka
func (c *Consumer) Consume(msgChan chan *ckafka.Message) error {
	// Infinite loop to consume messages
	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {
			msgChan <- msg
		} else {
			log.Printf("Error consuming message: %v (%v)\n", err, msg)
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() {
	c.consumer.Close()
}
