package sqlike

import (
	"context"

	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/types"
)

// Column : contains sql column information
type Column struct {
	// column name
	Name string

	// column position in sql database
	Position int

	// column data type with precision or size, eg. VARCHAR(20)
	Type string

	// column data type without precision and size, eg. VARCHAR
	DataType string

	// whether column is nullable or not
	IsNullable types.Boolean

	// default value of the column
	DefaultValue *string

	// text character set encoding
	Charset *string

	// text collation for sorting
	Collation *string

	// column comment
	Comment string

	// extra information
	Extra string
}

// ColumnView :
type ColumnView struct {
	tb *Table
}

// List : list all the column from current table
func (cv *ColumnView) List(ctx context.Context) ([]Column, error) {
	return cv.tb.ListColumns(ctx)
}

// Rename : rename your column name
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

// DropOne : drop column with name
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
