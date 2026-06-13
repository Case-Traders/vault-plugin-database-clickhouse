package clickhouse

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	F "github.com/IBM/fp-go/v2/function"
	IOR "github.com/IBM/fp-go/v2/idiomatic/ioresult"
	O "github.com/IBM/fp-go/v2/idiomatic/option"
	RES "github.com/IBM/fp-go/v2/idiomatic/result"
	P "github.com/IBM/fp-go/v2/predicate"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/template"
	"vault-plugin-database-clickhouse/internal/cluster"
	"vault-plugin-database-clickhouse/internal/deletepath"
	"vault-plugin-database-clickhouse/internal/stmt"
	"vault-plugin-database-clickhouse/internal/txexec"
	"vault-plugin-database-clickhouse/internal/user"
	"vault-plugin-database-clickhouse/internal/validate"
	"vault-plugin-database-clickhouse/internal/vars"
)

// New returns a Vault database plugin backed by ClickHouse.
func New() (interface{}, error) {
	db := new()
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
	return dbType, nil
}

func new() *Clickhouse {
	connProducer := &connutil.SQLConnectionProducer{}
	connProducer.Type = TypeName
	return &Clickhouse{SQLConnectionProducer: connProducer}
}

// Clickhouse implements dbplugin.Database for ClickHouse.
type Clickhouse struct {
	*connutil.SQLConnectionProducer

	usernameProducer template.StringTemplate
	clusterConfig    string
}

// Initialize parses config, opens the admin connection, and stores cluster and username settings.
func (p *Clickhouse) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	return RES.Chain(func(cfg map[string]interface{}) (dbplugin.InitializeResponse, error) {
		newConf, err := p.SQLConnectionProducer.Init(ctx, cfg, req.VerifyConnection)
		if err != nil {
			return dbplugin.InitializeResponse{}, err
		}

		up, err := setupUsernameTemplate(cfg)
		if err != nil {
			return dbplugin.InitializeResponse{}, err
		}
		p.usernameProducer = up

		clusterCfg, err := clusterFromConfig(cfg)
		if err != nil {
			return dbplugin.InitializeResponse{}, fmt.Errorf("retrieve %s: %w", ConfigKeyCluster, err)
		}
		p.clusterConfig = clusterCfg

		return dbplugin.InitializeResponse{Config: newConf}, nil
	})(req.Config, nil)
}

// Type reports the plugin name registered with Vault.
func (p *Clickhouse) Type() (string, error) {
	return TypeName, nil
}

func (p *Clickhouse) connection(ctx context.Context) IOR.IOResult[*sql.DB] {
	return func() (*sql.DB, error) {
		db, err := p.Connection(ctx)
		if err != nil {
			return nil, err
		}
		return db.(*sql.DB), nil
	}
}

func (p *Clickhouse) resolveCluster(ctx context.Context, db *sql.DB) (string, error) {
	return cluster.Resolve(ctx, db, p.clusterConfig)
}

func whenPresent[T any](value *T, run func(*T) error) error {
	v, ok := O.FromNillable(value)
	if !ok {
		return nil
	}
	return run(v)
}

func requireUser(ctx context.Context, db *sql.DB, username string) IOR.IOResult[*sql.DB] {
	return func() (*sql.DB, error) {
		exists, err := user.Exists(ctx, db, username)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("user %q does not exist", username)
		}
		return db, nil
	}
}

func ioDone(f func() error) IOR.IOResult[struct{}] {
	return func() (struct{}, error) {
		return struct{}{}, f()
	}
}

func (p *Clickhouse) withConnection(ctx context.Context, fn func(*sql.DB) error) error {
	_, err := F.Pipe1(
		p.connection(ctx),
		IOR.Chain(func(db *sql.DB) IOR.IOResult[struct{}] {
			return ioDone(func() error { return fn(db) })
		}),
	)()
	return err
}

func (p *Clickhouse) withExistingUser(ctx context.Context, username string, fn func(*sql.DB, string) error) error {
	return p.withConnection(ctx, func(db *sql.DB) error {
		if _, err := requireUser(ctx, db, username)(); err != nil {
			return err
		}
		clusterName, err := p.resolveCluster(ctx, db)
		if err != nil {
			return err
		}
		return fn(db, clusterName)
	})
}

// UpdateUser changes password and/or expiration for an existing ClickHouse user.
func (p *Clickhouse) UpdateUser(ctx context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	hasPassword := req.Password != nil
	hasExpiration := req.Expiration != nil
	if err := validate.UpdateUser(req.Username, hasPassword, hasExpiration); err != nil {
		return dbplugin.UpdateUserResponse{}, err
	}

	p.Lock()
	defer p.Unlock()

	return dbplugin.UpdateUserResponse{}, joinErrors([]error{
		whenPresent(req.Password, func(cp *dbplugin.ChangePassword) error {
			return p.changeUserPassword(ctx, req.Username, cp)
		}),
		whenPresent(req.Expiration, func(ce *dbplugin.ChangeExpiration) error {
			return p.changeUserExpiration(ctx, req.Username, ce)
		}),
	})
}

func joinErrors(errs []error) error {
	var merr *multierror.Error
	for _, e := range errs {
		if e != nil {
			merr = multierror.Append(merr, e)
		}
	}
	return merr.ErrorOrNil()
}

func (p *Clickhouse) changeUserPassword(ctx context.Context, username string, changePass *dbplugin.ChangePassword) error {
	password, err := RES.FromPredicate(P.IsNonZero[string](), func(string) error {
		return fmt.Errorf("missing password")
	})(changePass.NewPassword)
	if err != nil {
		return err
	}
	roleStmts := stmt.StatementsOrDefault(changePass.Statements.Commands, DefaultChangePasswordStatement)
	return p.withExistingUser(ctx, username, func(db *sql.DB, clusterName string) error {
		return p.runStatements(ctx, db, roleStmts, vars.ForUpdatePassword(username, password, clusterName))
	})
}

func (p *Clickhouse) changeUserExpiration(ctx context.Context, username string, changeExp *dbplugin.ChangeExpiration) error {
	roleStmts := stmt.StatementsOrDefault(changeExp.Statements.Commands, DefaultChangeExpirationStatement)
	return p.withExistingUser(ctx, username, func(db *sql.DB, clusterName string) error {
		tmpl := vars.ForUpdateExpiration(username, formatExpiration(changeExp.NewExpiration), clusterName)
		return p.runStatements(ctx, db, roleStmts, tmpl)
	})
}

// NewUser runs creation statements and returns the generated Vault username.
func (p *Clickhouse) NewUser(ctx context.Context, req dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	if err := validate.CreationStatements(len(req.Statements.Commands)); err != nil {
		return dbplugin.NewUserResponse{}, dbutil.ErrEmptyCreationStatement
	}

	p.Lock()
	defer p.Unlock()

	username, err := p.usernameProducer.Generate(req.UsernameConfig)
	if err != nil {
		return dbplugin.NewUserResponse{}, err
	}

	if err := p.withConnection(ctx, func(db *sql.DB) error {
		clusterName, err := p.resolveCluster(ctx, db)
		if err != nil {
			return err
		}
		tmpl := vars.ForNewUser(username, req.Password, formatExpiration(req.Expiration), clusterName)
		return p.runStatements(ctx, db, req.Statements.Commands, tmpl)
	}); err != nil {
		return dbplugin.NewUserResponse{}, err
	}
	return dbplugin.NewUserResponse{Username: username}, nil
}

// DeleteUser drops a user with custom revocation statements or the built-in DROP USER path.
func (p *Clickhouse) DeleteUser(ctx context.Context, req dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	p.Lock()
	defer p.Unlock()

	err := p.withConnection(ctx, func(db *sql.DB) error {
		if deletepath.UseCustomRevocation(req.Statements.Commands) {
			return p.customDeleteUser(ctx, req.Username, req.Statements.Commands)
		}
		return p.defaultDeleteUser(ctx, req.Username)
	})

	return dbplugin.DeleteUserResponse{}, err
}

func (p *Clickhouse) customDeleteUser(ctx context.Context, username string, revocationStmts []string) error {
	return p.withConnection(ctx, func(db *sql.DB) error {
		return p.runStatements(ctx, db, revocationStmts, vars.ForDeleteUser(username))
	})
}

func (p *Clickhouse) defaultDeleteUser(ctx context.Context, username string) error {
	return p.withConnection(ctx, func(db *sql.DB) error {
		exists, err := user.Exists(ctx, db, username)
		if err != nil {
			return err
		}
		if !exists {
			return nil
		}
		clusterName, err := p.resolveCluster(ctx, db)
		if err != nil {
			return err
		}
		_, err = db.ExecContext(ctx, "DROP USER IF EXISTS $1 ON CLUSTER $2", username, clusterName)
		return err
	})
}

func (p *Clickhouse) runStatements(ctx context.Context, db *sql.DB, commands []string, tmpl vars.TemplateVars) error {
	_, err := IOR.Bracket(
		func() (*sql.Tx, error) {
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				return nil, fmt.Errorf("begin transaction: %w", err)
			}
			return tx, nil
		},
		func(tx *sql.Tx) IOR.IOResult[struct{}] {
			return func() (struct{}, error) {
				if err := txexec.Execute(ctx, tx, commands, tmpl.Map()); err != nil {
					return struct{}{}, err
				}
				if err := tx.Commit(); err != nil {
					return struct{}{}, err
				}
				return struct{}{}, nil
			}
		},
		func(_ struct{}, err error) func(*sql.Tx) IOR.IOResult[struct{}] {
			return func(tx *sql.Tx) IOR.IOResult[struct{}] {
				return func() (struct{}, error) {
					if err != nil {
						_ = tx.Rollback()
					}
					return struct{}{}, nil
				}
			}
		},
	)()
	return err
}

func (p *Clickhouse) secretValues() map[string]string {
	return map[string]string{
		p.Password: "[password]",
	}
}
