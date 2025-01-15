package repo

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/NovruzovE/application-design-test/internal/core/entity"
	"github.com/NovruzovE/application-design-test/internal/core/usecase/order"
)

type RoomAvailabilityInMemRepo struct {
	store []*entity.RoomAvailability
	log   *slog.Logger
	mu    sync.RWMutex
}

var _ order.RoomAvailabilityRepository = (*RoomAvailabilityInMemRepo)(nil)

func NewRoomAvailabilityInMemRepo(log *slog.Logger) *RoomAvailabilityInMemRepo {
	return &RoomAvailabilityInMemRepo{
		store: []*entity.RoomAvailability{},
		log:   log,
	}
}

func (r *RoomAvailabilityInMemRepo) GetRoomAvailability(_ context.Context, order entity.Order) ([]*entity.RoomAvailability, error) {
	var result []*entity.RoomAvailability
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, ra := range r.store {
		if ra.HotelID == order.HotelID && ra.RoomID == order.RoomID && (ra.Date.After(order.From) || ra.Date.Equal(order.From)) && ra.Date.Before(order.To) && ra.Quota > 0 {
			result = append(result, ra)
		}
	}
	return result, nil
}

func (r *RoomAvailabilityInMemRepo) UpdateRoomAvailability(_ context.Context, availability []*entity.RoomAvailability) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, ra := range availability {
		for i, raStore := range r.store {
			if ra.HotelID == raStore.HotelID && ra.RoomID == raStore.RoomID && ra.Date == raStore.Date {
				r.store[i] = ra
			}
		}
	}
	return nil
}

func (r *RoomAvailabilityInMemRepo) PrepareRepo() {
	r.log.Debug("Preparing in-memory room availability repo")
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store = []*entity.RoomAvailability{
		{"reddison", "lux", date(2024, 1, 1), 1}, // nolint:exhaustivestruct
		{"reddison", "lux", date(2024, 1, 2), 1},
		{"reddison", "lux", date(2024, 1, 3), 1},
		{"reddison", "lux", date(2024, 1, 4), 1},
		{"reddison", "lux", date(2024, 1, 5), 0},
	}
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
