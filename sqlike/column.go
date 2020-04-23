package sqlike

import (
	"context"

	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/types"
)

// Column :
type Column struct {
	Name         string
	Position     int
	Type         string
	DataType     string
	IsNullable   types.Boolean
	DefaultValue *string
	Charset      *string
	Collation    *string
	Extra        string
}

// ColumnView :
type ColumnView struct {
	tb *Table
}

// List :
func (cv *ColumnView) List(ctx context.Context) ([]Column, error) {
	return cv.tb.ListColumns(ctx)
}

// Rename :
func (cv *ColumnView) Rename(ctx context.Context, oldColName, newColName string) error {
	stmt := sqlstmt.AcquireStmt(cv.tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	cv.tb.dialect.RenameColumn(stmt, cv.tb.dbName, cv.tb.name, oldColName, newColName)
	_, err := sqldriver.Execute(
		ctx,
		cv.tb.driver,
		stmt,
		cv.tb.logger,
	)
	return err
}

// DropOne :
func (cv *ColumnView) DropOne(ctx context.Context, name string) error {
	stmt := sqlstmt.AcquireStmt(cv.tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	cv.tb.dialect.DropColumn(stmt, cv.tb.dbName, cv.tb.name, name)
	_, err := sqldriver.Execute(
		ctx,
		cv.tb.driver,
		stmt,
		cv.tb.logger,
	)
	return err
}
