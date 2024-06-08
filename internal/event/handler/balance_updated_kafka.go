package handler

import (
	"fmt"
	"sync"

	"github.com/GabrielBrotas/eda-events/pkg/events"
	"github.com/GabrielBrotas/eda-events/pkg/kafka"
)

type BalanceUpdatedKafkaHandler struct {
	kafka *kafka.Producer
}

func NewBalanceUpdatedKafkaHandler(kafka *kafka.Producer) *BalanceUpdatedKafkaHandler {
	return &BalanceUpdatedKafkaHandler{kafka: kafka}
}

func (h *BalanceUpdatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.kafka.Publish("balances", nil, message)
	fmt.Printf("Event %s published to Kafka\n", message.GetName())
}
