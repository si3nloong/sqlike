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
	client   *Client
	context  context.Context
	driver   *sql.Tx
	dialect  sqldialect.Dialect
	registry *codec.Registry
	logger   logs.Logger
}

// SessionContext :
type SessionContext interface {
	Table(name string) *Table
}

// Table :
func (tx *Transaction) Table(name string) *Table {
	return &Table{
		dbName:   tx.dbName,
		name:     name,
		pk:       tx.pk,
		client:   tx.client,
		driver:   tx.driver,
		dialect:  tx.dialect,
		registry: tx.registry,
		logger:   tx.logger,
	}
	// return &Session{
	// 	dbName:   tx.dbName,
	// 	table:    name,
	// 	pk:       tx.pk,
	// 	tx:       tx,
	// 	registry: tx.registry,
	// }
}

// RollbackTransaction :
func (tx *Transaction) RollbackTransaction() error {
	return tx.driver.Rollback()
}

// CommitTransaction :
func (tx *Transaction) CommitTransaction() error {
	return tx.driver.Commit()
}
