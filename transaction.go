package sqlike

import (
	"context"
	"database/sql"

	sqlx "github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/sql/dialect"
)

var txnCtxKey struct{}

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

	logger sqlx.Logger
}

// RollbackTransaction : Rollback aborts the transaction.
func (tx *Transaction) RollbackTransaction() error {
	return tx.driver.Rollback()
}

// CommitTransaction : Commit commits the transaction.
func (tx *Transaction) CommitTransaction() error {
	return tx.driver.Commit()
}

// Value :
func (tx *Transaction) Value(key interface{}) interface{} {
	if key == &txnCtxKey {
		return tx
	}
	return tx.Context.Value(key)
}
