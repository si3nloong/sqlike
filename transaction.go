package sqlike

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlike/v2/db"
)

var txnCtxKey struct{}

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
	dialect db.Dialect

	// logger
	logger db.Logger
}

// Rollback : Rollback aborts the transaction.
func (tx *Transaction) Rollback() error {
	return tx.driver.Rollback()
}

// Commit : Commit commits the transaction.
func (tx *Transaction) Commit() error {
	return tx.driver.Commit()
}

// Value :
func (tx *Transaction) Value(key any) any {
	if key == &txnCtxKey {
		return tx
	}
	return tx.Context.Value(key)
}
