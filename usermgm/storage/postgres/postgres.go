package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"
)

const (
	pqUnique      = "23505"
	usrInvCodeDup = "user_account_invite_code_key"
)

type Storage struct {
	db     *sqlx.DB
	logger *logrus.Entry
}

// NewStorage returns a new Storage from the provides psql databse string
func NewStorage(dbstring string, logger *logrus.Entry) (*Storage, error) {
	db, err := sqlx.Connect("postgres", dbstring)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to postgres '%s'", dbstring)
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)
	return &Storage{db: db, logger: logger}, nil
}

func NewStorageDB(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

// RunMigration runs the migrations in the dir using the goose package
func (s *Storage) RunMigration(dir string) error {
	return goose.Run("up", s.db.DB, dir)
}

// NewTestStorage returns a Storage that uses an isolated database for testing purposes
// and a teardown function
func NewTestStorage(dbstring string, migrationDir string) (*Storage, func()) {
	db, teardown := MustNewDevelopmentDB(dbstring, migrationDir)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(time.Hour)

	return NewStorageDB(db), teardown
}

type pgTx struct{}

type tx struct {
	*sqlx.Tx
	committed *bool
}

func (s *Storage) NewTransacton(ctx context.Context) (context.Context, error) {
	t, err := s.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, pgTx{}, &tx{
		Tx:        t,
		committed: new(bool),
	}), nil
}

func (s *Storage) Commit(ctx context.Context) error {
	t := getTx(ctx)
	if t == nil {
		return fmt.Errorf("not a transaction")
	}
	if *t.committed {
		return nil
	}
	if err := t.Commit(); err != nil {
		return err
	}
	*t.committed = true
	return nil
}

func (s *Storage) Rollback(ctx context.Context) error {
	t, ok := ctx.Value(pgTx{}).(*tx)
	if !ok {
		return fmt.Errorf("not a transaction")
	}
	if *t.committed {
		return nil
	}
	return t.Rollback()
}

func getTx(ctx context.Context) *tx {
	if t, ok := ctx.Value(pgTx{}).(*tx); ok && !*t.committed {
		return t
	}
	return nil
}

// PrepareNamed prepares a named query in the current transaction (if begun) or with the db.
func (s *Storage) prepareNamed(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	if tx := getTx(ctx); tx != nil {
		return tx.PrepareNamedContext(ctx, query)
	}
	return s.db.PrepareNamedContext(ctx, query)
}

// queryer returns the transaction or database query executor.
func (s *Storage) queryer(ctx context.Context) sqlx.QueryerContext {
	if tx := getTx(ctx); tx != nil {
		return tx
	}
	return s.db
}

// execer returns the transaction or database exec executor.
func (s *Storage) execer(ctx context.Context) sqlx.ExtContext {
	if tx := getTx(ctx); tx != nil {
		return tx
	}
	return s.db
}
