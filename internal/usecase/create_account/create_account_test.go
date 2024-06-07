package create_account

import (
	"context"
	"testing"

	"github.com/GabrielBrotas/eda-events/internal/entity"
	"github.com/GabrielBrotas/eda-events/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j")
	uowMock := &mocks.UowMock{}
	uowMock.On("Do", mock.Anything, mock.Anything).Return(nil)

	// mock account gateway Save method
	uc := NewCreateAccountUseCase(uowMock)
	inputDto := CreateAccountInputDTO{
		ClientID: client.ID,
	}
	output, err := uc.Execute(context.Background(), inputDto)
	assert.Nil(t, err)
	assert.NotNil(t, output.ID)
	uowMock.AssertExpectations(t)
	uowMock.AssertNumberOfCalls(t, "Do", 1)
}
