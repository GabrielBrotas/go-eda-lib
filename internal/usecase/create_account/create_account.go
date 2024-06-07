package create_account

import (
	"context"

	"github.com/GabrielBrotas/eda-events/internal/entity"
	"github.com/GabrielBrotas/eda-events/internal/gateway"
	"github.com/GabrielBrotas/eda-events/pkg/uow"
)

type CreateAccountInputDTO struct {
	ClientID string `json:"client_id"`
}

type CreateAccountOutputDTO struct {
	ID string
}

type CreateAccountUseCase struct {
	Uow uow.UowInterface
}

func NewCreateAccountUseCase(uow uow.UowInterface) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		Uow: uow,
	}
}

func (uc *CreateAccountUseCase) Execute(ctx context.Context, input CreateAccountInputDTO) (*CreateAccountOutputDTO, error) {
	var accountID string
	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		clientRepository := uc.getClientRepository(ctx)
		accountRepository := uc.getAccountRepository(ctx)

		client, err := clientRepository.Get(input.ClientID)
		if err != nil {
			return err
		}

		account := entity.NewAccount(client)
		accountID = account.ID

		return accountRepository.Save(account)
	})

	if err != nil {
		return nil, err
	}

	output := &CreateAccountOutputDTO{
		ID: accountID,
	}

	return output, nil
}

func (uc *CreateAccountUseCase) getClientRepository(ctx context.Context) gateway.ClientGateway {
	repo, err := uc.Uow.GetRepository(ctx, "ClientDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.ClientGateway)
}

func (uc *CreateAccountUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}
