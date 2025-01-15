package order

import (
	"context"
	"github.com/NovruzovE/application-design-test/internal/core/entity"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order entity.Order) error
}

type RoomAvailabilityRepository interface {
	GetRoomAvailability(ctx context.Context, order entity.Order) ([]*entity.RoomAvailability, error)
	UpdateRoomAvailability(ctx context.Context, availability []*entity.RoomAvailability) error
}

type TransactionManager interface {
	Do(context.Context, func(context.Context) error) error
}
