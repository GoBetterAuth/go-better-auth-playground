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

	ID        string    `json:"id" bun:",pk"`
	EventType string    `json:"event_type"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
}
