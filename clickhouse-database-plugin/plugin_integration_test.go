//go:build integration

package main

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"vault-plugin-database-clickhouse/testutil"
)

// TestHelperProcess serves the database plugin when re-execed by go-plugin tests.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	if err := Run(); err != nil {
		t.Fatal(err)
	}
}

func pluginCmd(t *testing.T) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "-test.count=1", "--")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

func newRPCDatabase(t *testing.T, dsn string, extra map[string]interface{}) dbplugin.Database {
	t.Helper()

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  dbplugin.HandshakeConfig,
		VersionedPlugins: dbplugin.PluginSets,
		Cmd:              pluginCmd(t),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           hclog.NewNullLogger(),
	})
	t.Cleanup(func() { client.Kill() })

	rpc, err := client.Client()
	if err != nil {
		t.Fatalf("plugin client: %v", err)
	}

	raw, err := rpc.Dispense("database")
	if err != nil {
		t.Fatalf("dispense database: %v", err)
	}

	db, ok := raw.(dbplugin.Database)
	if !ok {
		t.Fatalf("unexpected dispense type %T", raw)
	}
	t.Cleanup(func() { _ = db.Close() })

	cfg := map[string]interface{}{"connection_url": dsn}
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

func TestPluginRPC_Initialize(t *testing.T) {
	_, dsn := testutil.StartClickHouse(t)
	newRPCDatabase(t, dsn, nil)
}

func TestPluginRPC_NewUserUpdateDelete(t *testing.T) {
	_, dsn := testutil.StartClickHouse(t)
	db := newRPCDatabase(t, dsn, nil)
	ctx := context.Background()

	resp, err := db.NewUser(ctx, dbplugin.NewUserRequest{
		Statements: dbplugin.Statements{Commands: []string{
			"CREATE USER IF NOT EXISTS '{{name}}' IDENTIFIED WITH plaintext_password BY '{{password}}';",
		}},
		UsernameConfig: dbplugin.UsernameMetadata{DisplayName: "rpc", RoleName: "role"},
		Password:       "s3cret-pass",
		Expiration:     time.Now().Add(24 * time.Hour),
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
