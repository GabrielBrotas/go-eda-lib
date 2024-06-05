package create_transaction

import (
	"context"
	"testing"

	"github.com/GabrielBrotas/eda-events/internal/entity"
	"github.com/GabrielBrotas/eda-events/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	// mock client1 account gateway
	client1, _ := entity.NewClient("client1", "j@j.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", account1.ID).Return(account1, nil)

	// mock client2 account gateway
	client2, _ := entity.NewClient("client2", "j@j2.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	accountMock.On("FindByID", account2.ID).Return(account2, nil)

	// mock update balance From
	accountMock.On("UpdateBalance", account1).Return(nil)

	// mock update balance To
	accountMock.On("UpdateBalance", account2).Return(nil)

	// mock transaction gateway
	transactionMock := &mocks.TransactionGatewayMock{}
	transactionMock.On("Create", mock.Anything).Return(nil)

	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}

	ctx := context.Background()

	uc := NewCreateTransactionUseCase(accountMock, transactionMock)
	output, err := uc.Execute(ctx, inputDto)
	assert.Nil(t, err)
	assert.NotNil(t, output)

	accountMock.AssertExpectations(t)
	transactionMock.AssertExpectations(t)
	accountMock.AssertNumberOfCalls(t, "FindByID", 2)
	accountMock.AssertNumberOfCalls(t, "UpdateBalance", 2)
	transactionMock.AssertNumberOfCalls(t, "Create", 1)

	// validate balance
	assert.Equal(t, 900.0, account1.Balance)
	assert.Equal(t, 1100.0, account2.Balance)
}
