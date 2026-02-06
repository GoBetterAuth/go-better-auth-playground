-- +goose Up

CREATE TABLE IF NOT EXISTS log_entries (
  id BIGSERIAL PRIMARY KEY,
  event_type VARCHAR(32) NOT NULL,
  details TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS log_entries;
