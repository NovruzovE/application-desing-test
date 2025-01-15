package order_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/NovruzovE/application-design-test/internal/core/entity"
	orderUseCase "github.com/NovruzovE/application-design-test/internal/core/usecase/order"
)

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

type MockRoomAvailabilityRepository struct {
	mock.Mock
}

func (m *MockRoomAvailabilityRepository) GetRoomAvailability(ctx context.Context, order entity.Order) ([]*entity.RoomAvailability, error) {
	args := m.Called(ctx, order)
	return args.Get(0).([]*entity.RoomAvailability), args.Error(1)
}

func (m *MockRoomAvailabilityRepository) UpdateRoomAvailability(ctx context.Context, availability []*entity.RoomAvailability) error {
	args := m.Called(ctx, availability)
	return args.Error(0)
}

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) SaveOrder(ctx context.Context, order entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func TestCreateOrder_Success(t *testing.T) {
	ctx := context.Background()
	order := entity.Order{
		HotelID:   "test",
		RoomID:    "test",
		UserEmail: "test",
		From:      time.Now(),
		To:        time.Now().Add(48 * time.Hour),
	}

	mockTrm := new(MockTransactionManager)
	mockRoomAvailRepo := new(MockRoomAvailabilityRepository)
	mockOrderRepo := new(MockOrderRepository)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mockRoomAvailRepo.On("GetRoomAvailability", ctx, order).Return([]*entity.RoomAvailability{
		{Quota: 2},
		{Quota: 2},
	}, nil)
	mockRoomAvailRepo.On("UpdateRoomAvailability", ctx, mock.Anything).Return(nil)
	mockOrderRepo.On("SaveOrder", ctx, order).Return(nil)
	mockTrm.On("Do", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(ctx context.Context) error)
		fn(ctx)
	})

	useCase := orderUseCase.NewOrderUseCase(mockRoomAvailRepo, mockOrderRepo, mockTrm, logger)
	err := useCase.CreateOrder(ctx, order)

	assert.NoError(t, err)
	mockRoomAvailRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
	mockTrm.AssertExpectations(t)
}
