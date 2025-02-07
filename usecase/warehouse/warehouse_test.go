package warehouse

import (
	"errors"
	"testing"
	"warehouse-service/models/warehouse"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockWarehouseRepository struct {
	mock.Mock
}

func (m *MockWarehouseRepository) Insert(warehouse *warehouse.RegisterRequest) error {
	args := m.Called(warehouse)
	return args.Error(0)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockWarehouseRepository)
	warehouseUsecase := NewWarehouseUsecase(mockRepo)

	mockRequest := &warehouse.RegisterRequest{
		Name:    "Warehouse A",
		Address: "123 Storage St",
	}

	// Expect Insert to be called with mockRequest and return nil (success)
	mockRepo.On("Insert", mockRequest).Return(nil)

	err := warehouseUsecase.Register(mockRequest)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegister_Fail(t *testing.T) {
	mockRepo := new(MockWarehouseRepository)
	warehouseUsecase := NewWarehouseUsecase(mockRepo)

	mockRequest := &warehouse.RegisterRequest{
		Name:    "Warehouse B",
		Address: "456 Industrial Rd",
	}

	// Simulate a database error
	mockRepo.On("Insert", mockRequest).Return(errors.New("database error"))

	err := warehouseUsecase.Register(mockRequest)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}
