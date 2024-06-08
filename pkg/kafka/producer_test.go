package kafka

import (
	"testing"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"
)

func TestProducerPublish(t *testing.T) {
	type TransactionDtoOutput struct {
		ID           string `json:"id"`
		Status       string `json:"status"`
		ErrorMessage string `json:"error_message"`
	}

	expectedOutput := TransactionDtoOutput{
		ID:           "1",
		Status:       "rejected",
		ErrorMessage: "you dont have limit for this transaction",
	}

	configMap := ckafka.ConfigMap{
		"test.mock.num.brokers": 3,
	}
	producer, err := NewProducer(&configMap)
	assert.Nil(t, err)
	err = producer.Publish("test", []byte("1"), expectedOutput)
	assert.Nil(t, err)
}
