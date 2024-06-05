package gateway

import "github.com/GabrielBrotas/eda-events/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}
