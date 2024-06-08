package kafka

import (
	"encoding/json"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *ckafka.Producer
}

func NewProducer(configMap *ckafka.ConfigMap) (*Producer, error) {
	p, err := ckafka.NewProducer(configMap)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: p}, nil
}

// Publish publishes a message to Kafka
func (p *Producer) Publish(topic string, key []byte, msg interface{}) error {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          msgJson,
		Key:            key,
	}
	err = p.producer.Produce(message, nil)
	if err != nil {
		panic(err)
	}
	return nil
}

func (p *Producer) Close() {
	p.producer.Close()
}
