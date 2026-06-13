//go:build integration

package clickhouse

import (
	"context"
	"testing"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"vault-plugin-database-clickhouse/testutil"
)

func newPlugin(t *testing.T, dsn string, extra map[string]interface{}) dbplugin.Database {
	t.Helper()
	raw, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	db := raw.(dbplugin.Database)
	cfg := map[string]interface{}{
		"connection_url": dsn,
	}
	for k, v := range extra {
		cfg[k] = v
	}
	if _, err := db.Initialize(context.Background(), dbplugin.InitializeRequest{
		Config:           cfg,
		VerifyConnection: true,
	}); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	return db
}

func TestIntegration_Initialize(t *testing.T) {
	_, dsn := testutil.StartClickHouse(t)
	newPlugin(t, dsn, nil)
}

func TestIntegration_NewUserUpdateDelete(t *testing.T) {
	_, dsn := testutil.StartClickHouse(t)
	db := newPlugin(t, dsn, nil)
	ctx := context.Background()

	create := []string{
		"CREATE USER IF NOT EXISTS '{{name}}' IDENTIFIED WITH plaintext_password BY '{{password}}';",
	}
	resp, err := db.NewUser(ctx, dbplugin.NewUserRequest{
		Statements: dbplugin.Statements{Commands: create},
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "role",
		},
		Password:   "s3cret-pass",
		Expiration: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("NewUser: %v", err)
	}
	if resp.Username == "" {
		t.Fatal("expected generated username")
	}

	_, err = db.UpdateUser(ctx, dbplugin.UpdateUserRequest{
		Username: resp.Username,
		Password: &dbplugin.ChangePassword{
			NewPassword: "new-pass",
			Statements: dbplugin.Statements{Commands: []string{
				"ALTER USER '{{username}}' IDENTIFIED WITH plaintext_password BY '{{password}}';",
			}},
		},
	})
	if err != nil {
		t.Fatalf("UpdateUser password: %v", err)
	}

	exp := time.Now().Add(48 * time.Hour)
	_, err = db.UpdateUser(ctx, dbplugin.UpdateUserRequest{
		Username: resp.Username,
		Expiration: &dbplugin.ChangeExpiration{
			NewExpiration: exp,
			Statements: dbplugin.Statements{Commands: []string{
				"ALTER USER '{{username}}' VALID UNTIL '{{expiration}}';",
			}},
		},
	})
	if err != nil {
		t.Fatalf("UpdateUser expiration: %v", err)
	}

	_, err = db.DeleteUser(ctx, dbplugin.DeleteUserRequest{
		Username: resp.Username,
		Statements: dbplugin.Statements{Commands: []string{
			"DROP USER IF EXISTS '{{username}}';",
		}},
	})
	if err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
}

func TestIntegration_UpdateUser_missingUser(t *testing.T) {
	_, dsn := testutil.StartClickHouse(t)
	db := newPlugin(t, dsn, nil)
	ctx := context.Background()

	_, err := db.UpdateUser(ctx, dbplugin.UpdateUserRequest{
		Username: "no-such-vault-user",
		Password: &dbplugin.ChangePassword{NewPassword: "x"},
	})
	if err == nil {
		t.Fatal("expected error updating missing user")
	}
}

func TestIntegration_clusterConfigOverride(t *testing.T) {
	_, dsn := testutil.StartClickHouse(t)
	db := newPlugin(t, dsn, map[string]interface{}{
		"cluster": "default",
	})
	ctx := context.Background()

	_, err := db.NewUser(ctx, dbplugin.NewUserRequest{
		Statements: dbplugin.Statements{Commands: []string{
			"CREATE USER IF NOT EXISTS '{{name}}' IDENTIFIED WITH plaintext_password BY '{{password}}';",
		}},
		UsernameConfig: dbplugin.UsernameMetadata{DisplayName: "cfg", RoleName: "r"},
		Password:       "pw",
		Expiration:     time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("NewUser with cluster config: %v", err)
	}
}
