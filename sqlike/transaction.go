package sqlike

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlike/sql/codec"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sqlike/logs"
)

// Transaction :
type Transaction struct {
	dbName   string
	pk       string
	context  context.Context
	driver   *sql.Tx
	dialect  sqldialect.Dialect
	registry *codec.Registry
	logger   logs.Logger
}

// Table :
func (tx *Transaction) Table(name string) *Session {
	return &Session{
		dbName:   tx.dbName,
		table:    name,
		pk:       tx.pk,
		tx:       tx,
		registry: tx.registry,
	}
}

// RollbackTransaction :
func (tx *Transaction) RollbackTransaction() error {
	return tx.driver.Rollback()
}

// CommitTransaction :
func (tx *Transaction) CommitTransaction() error {
	return tx.driver.Commit()
}
