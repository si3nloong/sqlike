package sqlike

import (
	"context"

	"github.com/si3nloong/sqlike/sqlike/logs"
	sqldriver "github.com/si3nloong/sqlike/sqlike/sql/driver"
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
	CharSet      *string
	Collation    *string
	Extra        string
}

// ColumnView :
type ColumnView struct {
	tb     *Table
	logger logs.Logger
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
		cv.tb.dialect.DropColumn(cv.tb.name, name),
		cv.tb.logger,
	)
	return err
}
