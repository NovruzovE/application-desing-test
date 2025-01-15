package repo

import (
	"context"
	"log/slog"
	"sync"

	"github.com/NovruzovE/application-design-test/internal/core/entity"
	"github.com/NovruzovE/application-design-test/internal/core/usecase/order"
)

type OrderInMemRepo struct {
	store []*entity.Order
	log   *slog.Logger
	mu    sync.RWMutex
}

var _ order.OrderRepository = (*OrderInMemRepo)(nil)

func NewOrderInMemRepo(logger *slog.Logger) *OrderInMemRepo {
	return &OrderInMemRepo{
		store: []*entity.Order{},
		log:   logger,
	}
}

func (r *OrderInMemRepo) SaveOrder(_ context.Context, order entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store = append(r.store, &order)

	return nil
}
