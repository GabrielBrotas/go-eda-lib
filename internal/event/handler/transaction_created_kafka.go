package handler

import (
	"fmt"
	"sync"

	"github.com/GabrielBrotas/eda-events/pkg/events"
	"github.com/GabrielBrotas/eda-events/pkg/kafka"
)

type TransactionCreatedKafkaHandler struct {
	kafka *kafka.Producer
}

func NewTransactionCreatedKafkaHandler(kafka *kafka.Producer) *TransactionCreatedKafkaHandler {
	return &TransactionCreatedKafkaHandler{kafka: kafka}
}

func (h *TransactionCreatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.kafka.Publish("transactions", nil, message)
	fmt.Printf("Event %s published to Kafka\n", message.GetName())
}
