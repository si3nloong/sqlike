package sqlike

import (
	sqldriver "bitbucket.org/SianLoong/sqlike/sqlike/sql/driver"
	"bitbucket.org/SianLoong/sqlike/types"
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
	tb *Table
}

// List :
func (cv *ColumnView) List() ([]Column, error) {
	return cv.tb.ListColumns()
}

// // AddOne :
// func (cv *ColumnView) AddOne() error {
// 	return nil
// }

// DropOne :
func (cv *ColumnView) DropOne(name string) error {
	_, err := sqldriver.Execute(
		cv.tb.driver,
		cv.tb.dialect.DropColumn(cv.tb.name, name),
		cv.tb.logger,
	)
	return err
}
