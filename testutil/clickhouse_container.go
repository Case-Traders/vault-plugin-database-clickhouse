package testutil

import (
	"context"
	"testing"

	tcclickhouse "github.com/testcontainers/testcontainers-go/modules/clickhouse"
)

// StartClickHouse runs clickhouse/clickhouse-server:24.8-alpine and returns a native DSN.
func StartClickHouse(t *testing.T) (*tcclickhouse.ClickHouseContainer, string) {
	t.Helper()
	ctx := context.Background()
	container, err := tcclickhouse.Run(ctx,
		"clickhouse/clickhouse-server:24.8-alpine",
	)
	if err != nil {
		t.Fatalf("start clickhouse: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("terminate clickhouse: %v", err)
		}
	})
	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}
	return container, dsn
}
