package create_client

import (
	"context"
	"time"

	"github.com/GabrielBrotas/eda-events/internal/entity"
	"github.com/GabrielBrotas/eda-events/internal/gateway"
	"github.com/GabrielBrotas/eda-events/pkg/uow"
)

type CreateClientInputDTO struct {
	Name  string
	Email string
}

type CreateClientOutputDTO struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateClientUseCase struct {
	Uow uow.UowInterface
}

func NewCreateClientUseCase(uow uow.UowInterface) *CreateClientUseCase {
	return &CreateClientUseCase{
		Uow: uow,
	}
}

func (uc *CreateClientUseCase) Execute(ctx context.Context, input CreateClientInputDTO) (*CreateClientOutputDTO, error) {
	client, err := entity.NewClient(input.Name, input.Email)
	if err != nil {
		return nil, err
	}

	err = uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		clientRepository := uc.getClientRepository(ctx)
		return clientRepository.Save(client)
	})

	if err != nil {
		return nil, err
	}

	output := &CreateClientOutputDTO{
		ID:        client.ID,
		Name:      client.Name,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
	return output, nil
}

func (uc *CreateClientUseCase) getClientRepository(ctx context.Context) gateway.ClientGateway {
	repo, err := uc.Uow.GetRepository(ctx, "ClientDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.ClientGateway)
}
