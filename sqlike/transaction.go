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
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryStmt(query interface{}) (*Result, error)
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

// Prepare : PrepareContext creates a prepared statement for use within a transaction.
func (tx *Transaction) Prepare(query string) (*sql.Stmt, error) {
	return tx.driver.PrepareContext(tx, query)
}

// Exec : ExecContext executes a query that doesn't return rows. For example: an INSERT and UPDATE.
func (tx *Transaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.driver.ExecContext(tx, query, args...)
}

// Query : QueryContext executes a query that returns rows, typically a SELECT.
func (tx *Transaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.driver.QueryContext(tx, query, args...)
}

// QueryRow : QueryRowContext executes a query that is expected to return at most one row. QueryRowContext always returns a non-nil value. Errors are deferred until Row's Scan method is called.
func (tx *Transaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return tx.driver.QueryRowContext(tx, query, args...)
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

// QueryStmt : QueryStmt support complex and advance query statement, make sure you executes a query that returns rows, typically a SELECT.
func (tx *Transaction) QueryStmt(query interface{}) (*Result, error) {
	if query == nil {
		return nil, errors.New("empty query statement")
	}
	stmt := sqlstmt.AcquireStmt(tx.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	if err := tx.dialect.SelectStmt(stmt, query); err != nil {
		return nil, err
	}
	rows, err := driver.Query(
		tx,
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

// RollbackTransaction : Rollback aborts the transaction.
func (tx *Transaction) RollbackTransaction() error {
	return tx.driver.Rollback()
}

// CommitTransaction : Commit commits the transaction.
func (tx *Transaction) CommitTransaction() error {
	return tx.driver.Commit()
}
