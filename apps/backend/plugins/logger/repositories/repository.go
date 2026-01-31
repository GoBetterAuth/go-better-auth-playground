package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/GoBetterAuth/go-better-auth-playground/plugins/logger/types"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// BunLoggerRepository implements Repository
type BunLoggerRepository struct {
	db bun.IDB
}

// NewBunLoggerRepository creates a new GORM-based repository
func NewBunLoggerRepository(db bun.IDB) *BunLoggerRepository {
	return &BunLoggerRepository{db: db}
}

// Create saves a new log entry to the database
func (r *BunLoggerRepository) Create(ctx context.Context, entry *types.LogEntry) error {
	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}

	if _, err := r.db.NewInsert().Model(entry).Exec(ctx); err != nil {
		return fmt.Errorf("failed to create log entry: %w", err)
	}
	return nil
}

// GetByID retrieves a log entry by ID
func (r *BunLoggerRepository) GetByID(ctx context.Context, id string) (*types.LogEntry, error) {
	var entry types.LogEntry
	if err := r.db.NewSelect().Model(&entry).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to get log entry: %w", err)
	}
	if entry.ID == "" {
		return nil, nil
	}
	return &entry, nil
}

// GetAll retrieves all log entries
func (r *BunLoggerRepository) GetAll(ctx context.Context) ([]types.LogEntry, error) {
	var entries []types.LogEntry
	if err := r.db.NewSelect().Model(&entries).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to get log entries: %w", err)
	}
	return entries, nil
}

// Delete removes a log entry by ID
func (r *BunLoggerRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.db.NewDelete().Model(&types.LogEntry{}).Where("id = ?", id).Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete log entry: %w", err)
	}
	return nil
}

// Count returns the total number of log entries
func (r *BunLoggerRepository) Count(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().Model(&types.LogEntry{}).Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count log entries: %w", err)
	}
	return count, nil
}

// Close closes the repository
func (r *BunLoggerRepository) Close() error {
	return nil
}
