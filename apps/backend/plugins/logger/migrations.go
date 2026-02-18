package logger

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/GoBetterAuth/go-better-auth/v2/migrations"
)

func loggerMigrations(provider string) []migrations.Migration {
	return migrations.ForProvider(provider, migrations.ProviderVariants{
		"sqlite": func() []migrations.Migration {
			return []migrations.Migration{loggerSQLiteInitial()}
		},
		"postgres": func() []migrations.Migration {
			return []migrations.Migration{loggerPostgresInitial()}
		},
		"mysql": func() []migrations.Migration {
			return []migrations.Migration{loggerMySQLInitial()}
		},
	})
}

func loggerSQLiteInitial() migrations.Migration {
	return migrations.Migration{
		Version: "20260201000000_logger_initial",
		Up: func(ctx context.Context, tx bun.Tx) error {
			return migrations.ExecStatements(
				ctx,
				tx,
				`CREATE TABLE IF NOT EXISTS log_entries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_type VARCHAR(32) NOT NULL,
  details TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);`,
			)
		},
		Down: func(ctx context.Context, tx bun.Tx) error {
			return migrations.ExecStatements(
				ctx,
				tx,
				`DROP TABLE IF EXISTS log_entries;`,
			)
		},
	}
}

func loggerPostgresInitial() migrations.Migration {
	return migrations.Migration{
		Version: "20260201000000_logger_initial",
		Up: func(ctx context.Context, tx bun.Tx) error {
			return migrations.ExecStatements(
				ctx,
				tx,
				`CREATE TABLE IF NOT EXISTS log_entries (
  id BIGSERIAL PRIMARY KEY,
  event_type VARCHAR(32) NOT NULL,
  details TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);`,
			)
		},
		Down: func(ctx context.Context, tx bun.Tx) error {
			return migrations.ExecStatements(
				ctx,
				tx,
				`DROP TABLE IF EXISTS log_entries;`,
			)
		},
	}
}

func loggerMySQLInitial() migrations.Migration {
	return migrations.Migration{
		Version: "20260201000000_logger_initial",
		Up: func(ctx context.Context, tx bun.Tx) error {
			return migrations.ExecStatements(
				ctx,
				tx,
				`CREATE TABLE IF NOT EXISTS log_entries (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  event_type VARCHAR(32) NOT NULL,
  details TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
			)
		},
		Down: func(ctx context.Context, tx bun.Tx) error {
			return migrations.ExecStatements(
				ctx,
				tx,
				`DROP TABLE IF EXISTS log_entries;`,
			)
		},
	}
}
