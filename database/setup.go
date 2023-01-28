package database

import (
	"context"
	"database/sql"
)

func RunMigrations(ctx context.Context, sqlDB *sql.DB) error {
	_, err := sqlDB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS tokens (
			client_id VARCHAR(255) NOT NULL,
			access_token VARCHAR(255) NOT NULL,
			refresh_token VARCHAR(255) NOT NULL,
			scopes VARCHAR(200),
			blocked TINYINT(1),
			expires_at DATETIME NOT NULL,
			refresh_expires_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			UNIQUE(client_id, access_token),
			UNIQUE(access_token, refresh_token)
		);

		CREATE INDEX IF NOT EXISTS idx_access_tokens ON tokens(access_token);
		CREATE INDEX IF NOT EXISTS idx_access_refresh_token ON tokens(access_token, refresh_token);`)

	return err
}

func RevertMigration(ctx context.Context, sqlDB *sql.DB) error {
	_, err := sqlDB.ExecContext(ctx, `
		DROP INDEX IF EXISTS idx_access_tokens;
		DROP INDEX IF EXISTS idx_access_refresh_token
		DROP TABLE IF EXISTS tokens;
	`)
	return err
}
