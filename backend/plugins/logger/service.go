package logger

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LogEntry struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	EventType string    `json:"event_type"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

type LoggerService struct {
	db *gorm.DB
}

// NewLoggerService initializes the service
func NewLoggerService(db *gorm.DB) *LoggerService {
	return &LoggerService{db: db}
}

// Log creates a new entry in the database
func (s *LoggerService) Log(eventType string, details string) error {
	entry := LogEntry{
		ID:        uuid.NewString(),
		EventType: eventType,
		Details:   details,
		CreatedAt: time.Now().UTC(),
	}

	return s.db.Create(&entry).Error
}
