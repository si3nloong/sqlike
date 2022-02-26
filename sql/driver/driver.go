package driver

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql/charset"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
)

// Info :
type Info interface {
	DriverName() string
	Charset() charset.Code
	Collate() string
}

// Queryer :
type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Driver :
type Driver interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Execute :
func Execute(ctx context.Context, driver Driver, stmt *sqlstmt.Statement, logger db.Logger) (result sql.Result, err error) {
	if logger != nil {
		stmt.StartTimer()
		defer func() {
			stmt.StopTimer()
			logger.Debug(stmt)
		}()
	}
	result, err = driver.ExecContext(ctx, stmt.String(), stmt.Args()...)
	return
}

// Query :
func Query(ctx context.Context, driver Driver, stmt *sqlstmt.Statement, logger db.Logger) (rows *sql.Rows, err error) {
	if logger != nil {
		stmt.StartTimer()
		defer func() {
			stmt.StopTimer()
			logger.Debug(stmt)
		}()
	}
	rows, err = driver.QueryContext(ctx, stmt.String(), stmt.Args()...)
	return
}

// QueryRowContext :
func QueryRowContext(ctx context.Context, driver Driver, stmt *sqlstmt.Statement, logger db.Logger) (row *sql.Row) {
	if logger != nil {
		stmt.StartTimer()
		defer func() {
			stmt.StopTimer()
			logger.Debug(stmt)
		}()
	}
	row = driver.QueryRowContext(ctx, stmt.String(), stmt.Args()...)
	return
}
