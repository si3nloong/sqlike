package sqlike

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sqlike/logs"
)

// SessionContext :
type SessionContext interface {
	context.Context
}

// Transaction :
type Transaction struct {
	// transaction context
	context.Context

	// database name
	dbName string

	// default primary key
	pk     string
	client *Client

	// transaction driver
	driver *sql.Tx

	// sql dialect
	dialect dialect.Dialect

	// codec
	codec  codec.Codecer
	logger logs.Logger
}

// RollbackTransaction : Rollback aborts the transaction.
func (tx *Transaction) RollbackTransaction() error {
	return tx.driver.Rollback()
}

// CommitTransaction : Commit commits the transaction.
func (tx *Transaction) CommitTransaction() error {
	return tx.driver.Commit()
}
