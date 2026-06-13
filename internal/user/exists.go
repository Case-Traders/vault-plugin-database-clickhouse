package user

import (
	"context"
	"database/sql"
	"fmt"
)

const existsQuery = `SELECT count() > 0 FROM system.users WHERE name = $1`

// Exists reports whether a ClickHouse user with the given name exists.
func Exists(ctx context.Context, db *sql.DB, username string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, existsQuery, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user exists: %w", err)
	}
	return exists, nil
}
