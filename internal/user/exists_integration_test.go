//go:build integration

package user_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"vault-plugin-database-clickhouse/internal/user"
	"vault-plugin-database-clickhouse/testutil"
)

func TestExists_defaultUser(t *testing.T) {
	_, dsn := testutil.StartClickHouse(t)
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	ok, err := user.Exists(context.Background(), db, "default")
	if err != nil {
		t.Fatalf("Exists: %v", err)
	}
	if !ok {
		t.Fatal("expected default user to exist")
	}
}

func TestExists_unknownUser(t *testing.T) {
	_, dsn := testutil.StartClickHouse(t)
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	ok, err := user.Exists(context.Background(), db, "vault-no-such-user")
	if err != nil {
		t.Fatalf("Exists: %v", err)
	}
	if ok {
		t.Fatal("expected unknown user to be absent")
	}
}
