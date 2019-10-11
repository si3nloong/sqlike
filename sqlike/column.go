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
func (cv *ColumnView) List() ([]Column, error) {
	return cv.tb.ListColumns()
}

// DropOne :
func (cv *ColumnView) DropOne(name string) error {
	_, err := sqldriver.Execute(
		context.Background(),
		cv.tb.driver,
		cv.tb.dialect.DropColumn(cv.tb.dbName, cv.tb.name, name),
		cv.tb.logger,
	)
	return err
}
