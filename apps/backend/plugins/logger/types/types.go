package types

import (
	"time"

	"github.com/uptrace/bun"
)

type LoggerPluginConfig struct {
	Enabled bool `json:"enabled" toml:"enabled"`
	// MaxLogCount is the maximum number of logs to keep before stopping
	MaxLogCount int `json:"max_log_count" toml:"max_log_count"`
}

// Validate validates the configuration
func (c *LoggerPluginConfig) Validate() error {
	if c.MaxLogCount <= 0 {
		c.MaxLogCount = 1000
	}
	return nil
}

type LogEntry struct {
	bun.BaseModel `bun:"table:log_entries"`

	ID        int64     `json:"id" bun:"column:id,pk,autoincrement"`
	EventType string    `json:"event_type" bun:"column:event_type"`
	Details   string    `json:"details" bun:"column:details"`
	CreatedAt time.Time `json:"created_at" bun:"column:created_at,default:current_timestamp"`
}
