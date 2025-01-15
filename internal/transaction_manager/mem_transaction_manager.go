package transaction_manager

import (
	"context"
	"log/slog"
	"sync"
)

type MemTransactionManager struct {
	mu  sync.RWMutex
	log *slog.Logger
}

func (mtrm *MemTransactionManager) Do(ctx context.Context, f func(context.Context) error) error {

	mtrm.log.Debug("Transaction started")
	defer mtrm.log.Debug("Transaction finished")

	return f(ctx)
}

func NewMemTransactionManager(log *slog.Logger) *MemTransactionManager {
	return &MemTransactionManager{
		log: log,
	}
}
