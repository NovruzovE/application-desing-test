package order

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/NovruzovE/application-design-test/internal/core/entity"
)

var ErrNoAvailableRooms = errors.New("no available rooms for selected dates")
var ErrCannotBookRoom = errors.New("cannot book room")

type OrderUseCase struct {
	log           *slog.Logger
	trm           TransactionManager
	roomAvailRepo RoomAvailabilityRepository
	orderRepo     OrderRepository
}

func NewOrderUseCase(
	roomAvailRepo RoomAvailabilityRepository,
	orderRepo OrderRepository,
	trm TransactionManager,
	logger *slog.Logger) *OrderUseCase {

	return &OrderUseCase{
		roomAvailRepo: roomAvailRepo,
		orderRepo:     orderRepo,
		trm:           trm,
		log:           logger,
	}
}

func (o *OrderUseCase) CreateOrder(ctx context.Context, order entity.Order) error {

	return o.trm.Do(ctx, func(ctx context.Context) error {
		roomAvailability, err := o.roomAvailRepo.GetRoomAvailability(ctx, order)
		if err != nil {
			o.log.Error("cannot get room availability", "error", err)
			return ErrCannotBookRoom
		}

		daysToBook := int(order.To.Truncate(24*time.Hour).Sub(order.From.Truncate(24*time.Hour)).Hours() / 24)

		if len(roomAvailability) < daysToBook {
			return ErrNoAvailableRooms
		}

		for _, availability := range roomAvailability {
			availability.Quota -= 1
		}

		err = o.roomAvailRepo.UpdateRoomAvailability(ctx, roomAvailability)
		if err != nil {
			o.log.Error("cannot update room availability", "error", err)
			return ErrCannotBookRoom
		}

		err = o.orderRepo.SaveOrder(ctx, order)
		if err != nil {
			o.log.Error("cannot save order", "error", err)
			return ErrCannotBookRoom
		}

		return nil
	})
}
