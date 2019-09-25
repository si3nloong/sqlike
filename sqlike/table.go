package sqlike

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/si3nloong/sqlike/core"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/codec"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sql/util"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/si3nloong/sqlike/sqlike/logs"
)

// ErrNoRecordAffected :
var ErrNoRecordAffected = errors.New("no record affected")

// ErrExpectedStruct :
var ErrExpectedStruct = errors.New("expected struct as a source")

// ErrEmptyFields :
var ErrEmptyFields = errors.New("empty fields")

// INSERT INTO Table (X1, X2) VALUES (?,?)
// - required >> table: string, columns: []string, arguments :[][]interface{}
// - options >> omitFields: []string, set upsert

// Table :
type Table struct {
	dbName   string
	name     string
	pk       string
	client   *Client
	driver   sqldriver.Driver
	dialect  sqldialect.Dialect
	registry *codec.Registry
	logger   logs.Logger
}

// Rename : rename the current table name to new table name
func (tb *Table) Rename(name string) error {
	_, err := sqldriver.Execute(
		context.Background(),
		tb.driver,
		tb.dialect.RenameTable(tb.dbName, tb.name, name),
		tb.logger,
	)
	return err
}

// Exists : this will return true when the table exists in the database
func (tb *Table) Exists() bool {
	var count int
	stmt := tb.dialect.HasTable(tb.dbName, tb.name)
	row := sqldriver.QueryRowContext(
		context.Background(),
		tb.driver,
		stmt,
		tb.logger,
	)
	row.Scan(&count)
	return count > 0
}

// Columns :
func (tb *Table) Columns() *ColumnView {
	return &ColumnView{tb: tb}
}

// ListColumns :
func (tb *Table) ListColumns() ([]Column, error) {
	stmt := tb.dialect.GetColumns(tb.dbName, tb.name)
	rows, err := sqldriver.Query(
		context.Background(),
		tb.driver,
		stmt,
		tb.logger,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make([]Column, 0)
	for i := 0; rows.Next(); i++ {
		col := Column{}

		if err := rows.Scan(
			&col.Position,
			&col.Name,
			&col.Type,
			&col.DefaultValue,
			&col.IsNullable,
			&col.DataType,
			&col.CharSet,
			&col.Collation,
			&col.Extra,
		); err != nil {
			return nil, err
		}

		col.Type = strings.ToUpper(col.Type)
		col.DataType = strings.ToUpper(col.DataType)

		columns = append(columns, col)
	}
	return columns, nil
}

// ListIndexes :
func (tb *Table) ListIndexes() ([]Index, error) {
	stmt := tb.dialect.GetIndexes(tb.dbName, tb.name)
	rows, err := sqldriver.Query(
		context.Background(),
		tb.driver,
		stmt,
		tb.logger,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	idxs := make([]Index, 0)
	for i := 0; rows.Next(); i++ {
		idx := Index{}

		if err := rows.Scan(
			&idx.Name,
			&idx.Type,
			&idx.IsVisible,
		); err != nil {
			return nil, err
		}

		idxs = append(idxs, idx)
	}
	return idxs, nil
}

// MustMigrate : this will ensure the migrate is complete, otherwise it will panic
func (tb Table) MustMigrate(entity interface{}) {
	err := tb.Migrate(entity)
	if err != nil {
		panic(err)
	}
}

// Migrate : migrate will create a new table follows by the definition of struct tag, alter when the table already exists
func (tb *Table) Migrate(entity interface{}) error {
	return tb.migrateOne(entity, false)
}

// UnsafeMigrate : unsafe migration will delete non-exist
// index and columns, beware when you use this
func (tb Table) UnsafeMigrate(entity interface{}) error {
	return tb.migrateOne(entity, true)
}

// MustUnsafeMigrate :
func (tb Table) MustUnsafeMigrate(entity interface{}) {
	err := tb.migrateOne(entity, true)
	if err != nil {
		panic(err)
	}
}

// Truncate :
func (tb *Table) Truncate() (err error) {
	_, err = sqldriver.Execute(
		context.Background(),
		tb.driver,
		tb.dialect.TruncateTable(tb.dbName, tb.name),
		tb.logger,
	)
	return
}

// DropIfExits : will drop the table only if it exists
func (tb Table) DropIfExits() (err error) {
	_, err = sqldriver.Execute(
		context.Background(),
		tb.driver,
		tb.dialect.DropTable(tb.dbName, tb.name, true),
		tb.logger,
	)
	return
}

// Drop : drop the table, but it might throw error when the table is not exists
func (tb Table) Drop() (err error) {
	_, err = sqldriver.Execute(
		context.Background(),
		tb.driver,
		tb.dialect.DropTable(tb.dbName, tb.name, false),
		tb.logger,
	)
	return
}

// Copy :
func (tb *Table) Copy(fields []string, act actions.CopyStatement) error {
	x := new(actions.CopyActions)
	if act != nil {
		*x = *(act.(*actions.CopyActions))
	}
	if x.Table == "" {
		return errors.New("sqlike: empty table name")
	}
	stmt, err := tb.dialect.Copy(tb.dbName, tb.name, fields, x)
	if err != nil {
		return err
	}
	_, err = sqldriver.Execute(
		context.Background(),
		tb.driver,
		stmt,
		tb.logger,
	)
	return err
}

// Indexes :
func (tb *Table) Indexes() *IndexView {
	return &IndexView{tb: tb}
}

// HasIndexByName :
func (tb *Table) HasIndexByName(name string) (bool, error) {
	return isIndexExists(
		tb.dbName,
		tb.name,
		name,
		tb.driver,
		tb.dialect,
		tb.logger,
	)
}

func (tb *Table) migrateOne(entity interface{}, unsafe bool) error {
	v := reflext.ValueOf(entity)
	if !v.IsValid() {
		return ErrInvalidInput
	}

	t := reflext.Deref(v.Type())
	if !reflext.IsKind(t, reflect.Struct) {
		return ErrExpectedStruct
	}

	cdc := core.DefaultMapper.CodecByType(t)
	_, fields := skipColumns(cdc.Properties, nil)
	if len(fields) < 1 {
		return ErrEmptyFields
	}

	if !tb.Exists() {
		return tb.createTable(fields)
	}

	columns, err := tb.ListColumns()
	if err != nil {
		return err
	}
	return tb.alterTable(fields, columns, nil, unsafe)
}

func (tb *Table) createTable(fields []*reflext.StructField) error {
	stmt, err := tb.dialect.CreateTable(tb.dbName, tb.name, tb.pk, fields)
	if err != nil {
		return err
	}
	if _, err := sqldriver.Execute(
		context.Background(),
		tb.driver,
		stmt,
		tb.logger,
	); err != nil {
		return err
	}
	return nil
}

func (tb *Table) alterTable(fields []*reflext.StructField, columns []Column, indexs []indexes.Index, unsafe bool) error {
	cols := make([]string, len(columns), len(columns))
	for i, col := range columns {
		cols[i] = col.Name
	}
	idxs := make([]string, len(indexs), len(indexs))
	for i, idx := range indexs {
		idxs[i] = idx.Name
	}
	stmt, err := tb.dialect.AlterTable(tb.dbName, tb.name, tb.pk, fields, cols, idxs, unsafe)
	if err != nil {
		return err
	}
	if _, err := sqldriver.Execute(
		context.Background(),
		tb.driver,
		stmt,
		tb.logger,
	); err != nil {
		return err
	}
	return nil
}

// we should skip virtual columns on insertion and migration
func skipColumns(sfs []*reflext.StructField, omits util.StringSlice) (columns []string, fields []*reflext.StructField) {
	fields = make([]*reflext.StructField, 0, len(sfs))
	length := len(omits)
	for _, sf := range sfs {
		if _, isOk := sf.Tag.LookUp("generated"); isOk {
			continue
		}
		if length > 0 && omits.IndexOf(sf.Path) > -1 {
			continue
		}
		columns = append(columns, sf.Path)
		fields = append(fields, sf)
	}
	return
}
