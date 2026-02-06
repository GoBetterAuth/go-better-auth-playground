package repositories

import (
	"context"

	"github.com/GoBetterAuth/go-better-auth-playground/plugins/logger/types"
)

// LoggerRepository defines the interface for log entry persistence
type LoggerRepository interface {
	Create(ctx context.Context, entry *types.LogEntry) error
	GetByID(ctx context.Context, id int64) (*types.LogEntry, error)
	GetAll(ctx context.Context) ([]types.LogEntry, error)
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int, error)
	Close() error
}
