package sqlike

import (
	"context"
	"database/sql"

	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sqlike/logs"
)

// Transaction :
type Transaction struct {
	dbName  string
	pk      string
	context context.Context
	driver  *sql.Tx
	dialect sqldialect.Dialect
	logger  logs.Logger
}

// Table :
func (tx *Transaction) Table(name string) *Session {
	return &Session{
		dbName: tx.dbName,
		table:  name,
		pk:     tx.pk,
		tx:     tx,
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
