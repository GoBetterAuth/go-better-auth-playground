package services

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/GoBetterAuth/go-better-auth-playground/plugins/logger/repositories"
	"github.com/GoBetterAuth/go-better-auth-playground/plugins/logger/types"
	"github.com/GoBetterAuth/go-better-auth/v2/models"
	"github.com/google/uuid"
)

// service implements the UseCase interface for logger operations
type service struct {
	repo     repositories.LoggerRepository
	logger   models.Logger
	config   types.LoggerPluginConfig
	logCount atomic.Int64
}

// NewService creates a new logger usecase implementation
func NewService(repo repositories.LoggerRepository, logger models.Logger, config types.LoggerPluginConfig) LoggerService {
	return &service{
		repo:   repo,
		logger: logger,
		config: config,
	}
}

// CreateLogEntry creates a new log entry with hooks
func (s *service) CreateLogEntry(ctx context.Context, eventType string, details string) (*types.LogEntry, error) {
	entry := &types.LogEntry{
		ID:        uuid.NewString(),
		EventType: eventType,
		Details:   details,
		CreatedAt: time.Now().UTC(),
	}

	// Create in database
	if err := s.repo.Create(ctx, entry); err != nil {
		s.logger.Error("failed to create log entry", "error", err)
		return entry, err
	}

	// Refresh from database
	retrievedEntry, err := s.repo.GetByID(ctx, entry.ID)
	if err != nil {
		s.logger.Error("failed to retrieve created log entry", "error", err)
		return entry, err
	}
	if retrievedEntry == nil {
		return entry, fmt.Errorf("created entry not found in database")
	}

	s.logCount.Add(1)
	return retrievedEntry, nil
}

// GetLogEntry retrieves a log entry by ID
func (s *service) GetLogEntry(ctx context.Context, id string) (*types.LogEntry, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAllLogs retrieves all log entries
func (s *service) GetAllLogs(ctx context.Context) ([]types.LogEntry, error) {
	return s.repo.GetAll(ctx)
}

// DeleteLogEntry deletes a log entry
func (s *service) DeleteLogEntry(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetLogCount returns the current number of logs
func (s *service) GetLogCount(ctx context.Context) (int64, error) {
	return s.logCount.Load(), nil
}

// HasReachedMaxLogs checks if the maximum log count has been reached
func (s *service) HasReachedMaxLogs(ctx context.Context) (bool, error) {
	return int(s.logCount.Load()) >= s.config.MaxLogCount, nil
}
