package txexec

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"vault-plugin-database-clickhouse/internal/stmt"
)

// Execute runs normalized DDL statements in order; stops at the first failure.
func Execute(ctx context.Context, tx *sql.Tx, commands []string, vars map[string]string) error {
	queries := stmt.NormalizeCommands(commands)
	return FirstError(queries, func(q string) error {
		if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, vars, q); err != nil {
			return fmt.Errorf("execute query %q: %w", q, err)
		}
		return nil
	})
}
