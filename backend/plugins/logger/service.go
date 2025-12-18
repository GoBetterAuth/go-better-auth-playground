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
	db            *gorm.DB
	databaseHooks *LogEntryDatabaseHooks
}

// NewLoggerService initializes the service
func NewLoggerService(db *gorm.DB, databaseHooks *LogEntryDatabaseHooks) *LoggerService {
	return &LoggerService{db: db, databaseHooks: databaseHooks}
}

// CreateLogEntry creates a new entry in the database
func (s *LoggerService) CreateLogEntry(eventType string, details string) (*LogEntry, error) {
	entry := LogEntry{
		ID:        uuid.NewString(),
		EventType: eventType,
		Details:   details,
		CreatedAt: time.Now().UTC(),
	}

	if s.databaseHooks != nil && s.databaseHooks.BeforeCreate != nil {
		if err := s.databaseHooks.BeforeCreate(&entry); err != nil {
			return &entry, err
		}
	}

	if err := s.db.Create(&entry).Error; err != nil {
		return &entry, err
	}

	var logEntryCreated LogEntry
	if err := s.db.First(&logEntryCreated, "id = ?", entry.ID).Error; err != nil {
		return &entry, err
	}

	if s.databaseHooks != nil && s.databaseHooks.AfterCreate != nil {
		if err := s.databaseHooks.AfterCreate(logEntryCreated); err != nil {
			return &logEntryCreated, err
		}
	}

	return &logEntryCreated, nil
}
