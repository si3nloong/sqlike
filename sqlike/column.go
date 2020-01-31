package sqlike

import (
	"context"

	sqldriver "github.com/si3nloong/sqlike/sql/driver"
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
	_, err := sqldriver.Execute(
		ctx,
		cv.tb.driver,
		cv.tb.dialect.RenameColumn(cv.tb.dbName, cv.tb.name, oldColName, newColName),
		cv.tb.logger,
	)
	return err
}

// DropOne :
func (cv *ColumnView) DropOne(ctx context.Context, name string) error {
	_, err := sqldriver.Execute(
		ctx,
		cv.tb.driver,
		cv.tb.dialect.DropColumn(cv.tb.dbName, cv.tb.name, name),
		cv.tb.logger,
	)
	return err
}
