package gateway

import "github.com/GabrielBrotas/eda-events/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
