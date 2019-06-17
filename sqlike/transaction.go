package sqlike

import (
	"database/sql"

	sqlcore "github.com/si3nloong/sqlike/sqlike/sql/core"
)

// Session :
type Session struct {
	table string
	tx    *Transaction
}

// Transaction :
type Transaction struct {
	driver  *sql.Tx
	dialect sqlcore.Dialect
	logger  Logger
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
