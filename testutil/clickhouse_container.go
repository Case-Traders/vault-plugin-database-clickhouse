package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	tcclickhouse "github.com/testcontainers/testcontainers-go/modules/clickhouse"
)

// StartClickHouse runs clickhouse/clickhouse-server:24.8-alpine and returns a native DSN.
func StartClickHouse(t *testing.T) (*tcclickhouse.ClickHouseContainer, string) {
	t.Helper()
	ctx := context.Background()
	container, err := tcclickhouse.Run(ctx,
		"clickhouse/clickhouse-server:24.8-alpine",
		useDefaultSuperuser(),
	)
	if err != nil {
		t.Fatalf("start clickhouse: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("terminate clickhouse: %v", err)
		}
	})
	host, err := container.ConnectionHost(ctx)
	if err != nil {
		t.Fatalf("connection host: %v", err)
	}
	dsn := fmt.Sprintf("clickhouse://default@%s/default", host)
	return container, dsn
}

// useDefaultSuperuser clears CLICKHOUSE_USER so the image default superuser stays.
func useDefaultSuperuser() testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		delete(req.Env, "CLICKHOUSE_USER")
		delete(req.Env, "CLICKHOUSE_PASSWORD")
		return nil
	}
}
