package sqlike

import (
	"context"
	"database/sql"
	"errors"

	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
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
	// transaction context
	context.Context

	// database name
	dbName string

	// default primary key
	pk      string
	client  *Client
	driver  *sql.Tx
	dialect dialect.Dialect
	codec   codec.Codecer
	logger  logs.Logger
}

// Table :
func (tx *Transaction) Table(name string) *Table {
	return &Table{
		dbName:  tx.dbName,
		name:    name,
		pk:      tx.pk,
		client:  tx.client,
		driver:  tx.driver,
		dialect: tx.dialect,
		codec:   tx.codec,
		logger:  tx.logger,
	}
}

// QueryStmt : support complex and advance query statement
func (tx *Transaction) QueryStmt(ctx context.Context, query interface{}) (*Result, error) {
	if query == nil {
		return nil, errors.New("empty query statement")
	}
	stmt := sqlstmt.AcquireStmt(tx.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	if err := tx.dialect.SelectStmt(stmt, query); err != nil {
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
	rslt.cache = tx.client.cache
	rslt.codec = tx.codec
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
