package sqlike

import (
	"context"
	"database/sql"
	"errors"

	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/logs"
)

// SessionContext :
type SessionContext interface {
	context.Context
	Table(name string) *Table
	QueryStmt(ctx context.Context, query interface{}) (*Result, error)
}

// Transaction :
type Transaction struct {
	context.Context
	dbName   string
	pk       string
	client   *Client
	driver   *sql.Tx
	dialect  dialect.Dialect
	registry *codec.Registry
	logger   logs.Logger
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
}

func (tx *Transaction) QueryStmt(ctx context.Context, query interface{}) (*Result, error) {
	if query == nil {
		return nil, errors.New("empty query statement")
	}
	stmt, err := tx.dialect.SelectStmt(query)
	if err != nil {
		return nil, err
	}
	rows, err := driver.Query(
		ctx,
		tx.driver,
		stmt,
		getLogger(tx.logger, true),
	)
	if err != nil {
		return nil, err
	}
	rslt := new(Result)
	rslt.registry = tx.registry
	rslt.rows = rows
	rslt.columns, rslt.err = rows.Columns()
	return rslt, rslt.err
}

// RollbackTransaction :
func (tx *Transaction) RollbackTransaction() error {
	return tx.driver.Rollback()
}

// CommitTransaction :
func (tx *Transaction) CommitTransaction() error {
	return tx.driver.Commit()
}
