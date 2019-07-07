package sqldriver

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlike/sqlike/logs"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
)

// Driver :
type Driver interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Execute :
func Execute(ctx context.Context, driver Driver, stmt *sqlstmt.Statement, logger logs.Logger) (result sql.Result, err error) {
	if logger != nil {
		stmt.StartTimer()
		defer func() {
			stmt.StopTimer()
			logger.Format(stmt)
		}()
	}
	result, err = driver.ExecContext(ctx, stmt.String(), stmt.Args()...)
	return
}

// Query :
func Query(ctx context.Context, driver Driver, stmt *sqlstmt.Statement, logger logs.Logger) (rows *sql.Rows, err error) {
	if logger != nil {
		stmt.StartTimer()
		defer func() {
			stmt.StopTimer()
			logger.Format(stmt)
		}()
	}
	rows, err = driver.QueryContext(ctx, stmt.String(), stmt.Args()...)
	return
}

// QueryRowContext :
func QueryRowContext(ctx context.Context, driver Driver, stmt *sqlstmt.Statement, logger logs.Logger) (row *sql.Row) {
	if logger != nil {
		stmt.StartTimer()
		defer func() {
			stmt.StopTimer()
			logger.Format(stmt)
		}()
	}
	row = driver.QueryRowContext(ctx, stmt.String(), stmt.Args()...)
	return
}
