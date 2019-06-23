package sqlike

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlike/sqlike/logs"
	sqlcore "github.com/si3nloong/sqlike/sqlike/sql/core"
)

// Transaction :
type Transaction struct {
	context context.Context
	driver  *sql.Tx
	dialect sqlcore.Dialect
	logger  logs.Logger
}

// Table :
func (tx *Transaction) Table(name string) *Session {
	return &Session{
		table: name,
		tx:    tx,
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
