package sqlike

import (
	"context"
	"database/sql"
	"time"

	"github.com/si3nloong/sqlike/sqlike/logs"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/sql/codec"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
)

// Database :
type Database struct {
	name     string
	pk       string
	client   *Client
	driver   sqldriver.Driver
	dialect  sqldialect.Dialect
	registry *codec.Registry
	logger   logs.Logger
}

// Table :
func (db *Database) Table(name string) *Table {
	return &Table{
		dbName:   db.name,
		name:     name,
		pk:       db.pk,
		client:   db.client,
		driver:   db.driver,
		dialect:  db.dialect,
		registry: db.registry,
		logger:   db.logger,
	}
}

// BeginTransaction :
func (db *Database) BeginTransaction() (*Transaction, error) {
	return db.beginTrans(context.Background(), nil)
}

func (db *Database) beginTrans(ctx context.Context, opt *sql.TxOptions) (*Transaction, error) {
	tx, err := db.client.BeginTx(ctx, opt)
	if err != nil {
		return nil, err
	}
	return &Transaction{
		pk:      db.pk,
		context: ctx,
		driver:  tx,
		dialect: db.dialect,
		logger:  db.logger,
	}, nil
}

type txCallback func(ctx SessionContext) error

// RunInTransaction :
func (db *Database) RunInTransaction(cb txCallback, opts ...*options.TransactionOptions) error {
	opt := new(options.TransactionOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	duration := 60 * time.Second
	if opt.Duration.Seconds() > 0 {
		duration = opt.Duration
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	tx, err := db.beginTrans(ctx, &sql.TxOptions{
		Isolation: opt.IsolationLevel,
		ReadOnly:  opt.ReadOnly,
	})
	if err != nil {
		return err
	}
	defer tx.RollbackTransaction()
	if err := cb(tx); err != nil {
		return err
	}
	return tx.CommitTransaction()
}
