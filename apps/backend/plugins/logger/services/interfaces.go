package services

import (
	"context"

	"github.com/GoBetterAuth/go-better-auth-playground/plugins/logger/types"
)

// UseCase defines the interface for logger operations
type LoggerService interface {
	CreateLogEntry(ctx context.Context, eventType string, details string) (*types.LogEntry, error)
	GetLogEntry(ctx context.Context, id string) (*types.LogEntry, error)
	GetAllLogs(ctx context.Context) ([]types.LogEntry, error)
	DeleteLogEntry(ctx context.Context, id string) error
	GetLogCount(ctx context.Context) (int64, error)
	HasReachedMaxLogs(ctx context.Context) (bool, error)
}
