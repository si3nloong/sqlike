package sqlike

import (
	"reflect"
	"strings"
	"github.com/si3nloong/sqlike/core"
	"github.com/si3nloong/sqlike/core/codec"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	sqlcore "github.com/si3nloong/sqlike/sqlike/sql/core"
	sqldriver "github.com/si3nloong/sqlike/sqlike/sql/driver"
	"golang.org/x/xerrors"
)

// ErrNoRecordAffected :
var ErrNoRecordAffected = xerrors.New("no record affected")

// ErrExpectedStruct :
var ErrExpectedStruct = xerrors.New("expected struct as a source")

// ErrEmptyFields :
var ErrEmptyFields = xerrors.New("empty fields")

// INSERT INTO Table (X1, X2) VALUES (?,?)
// - required >> table: string, columns: []string, arguments :[][]interface{}
// - options >> omitFields: []string, set upsert

// Table :
type Table struct {
	dbName   string
	name     string
	client   *Client
	driver   sqldriver.Driver
	dialect  sqlcore.Dialect
	registry *codec.Registry
	logger   Logger
}

// Exists :
func (tb *Table) Exists() bool {
	var count int
	stmt := tb.dialect.HasTable(tb.dbName, tb.name)
	row := sqldriver.QueryRow(
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

// MustMigrate :
func (tb Table) MustMigrate(entity interface{}) {
	err := tb.Migrate(entity)
	if err != nil {
		panic(err)
	}
}

// Migrate :
func (tb *Table) Migrate(entity interface{}) error {
	return tb.migrateOne(entity, false)
}

// UnsafeMigrate : unsafe migration will delete non-exist
// index and columns, beware to use this
func (tb Table) UnsafeMigrate(entity interface{}) error {
	return tb.migrateOne(entity, true)
}

// Truncate :
func (tb Table) Truncate() (err error) {
	_, err = sqldriver.Execute(
		tb.driver,
		tb.dialect.TruncateTable(tb.name),
		tb.logger,
	)
	return
}

// Drop :
func (tb Table) Drop() (err error) {
	_, err = sqldriver.Execute(
		tb.driver,
		tb.dialect.DropTable(tb.name, false),
		tb.logger,
	)
	return
}

// ReplaceInto :
func (tb Table) ReplaceInto(filter interface{}) error {
	// stmt, args, err := tb.dialect.ReplaceInto(tb.name, filter)
	// if err != nil {
	// 	return err
	// }
	// _, err = sqldriver.Execute(
	// 	tb.driver,
	// 	stmt,
	// 	args,
	// 	tb.logger,
	// )
	// return err
	return nil
}

// Indexes :
func (tb *Table) Indexes() *IndexView {
	return &IndexView{tb: tb}
}

func (tb *Table) migrateOne(entity interface{}, unsafe bool) error {
	t := reflext.Deref(reflect.TypeOf(entity))
	if !reflext.IsKind(t, reflect.Struct) {
		return ErrExpectedStruct
	}

	cdc := core.DefaultMapper.CodecByType(t)
	_, fields := skipGeneratedColumns(cdc.NameFields)
	if len(fields) < 1 {
		return ErrEmptyFields
	}

	if tb.Exists() {
		columns, err := tb.ListColumns()
		if err != nil {
			return err
		}
		return tb.alterTable(fields, columns, nil, unsafe)
	}
	return tb.createTable(fields)
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
	stmt, err := tb.dialect.AlterTable(tb.name, fields, cols, idxs, unsafe)
	if err != nil {
		return err
	}
	if _, err := sqldriver.Execute(
		tb.driver,
		stmt,
		tb.logger,
	); err != nil {
		return err
	}
	return nil
}

func (tb *Table) createTable(fields []*reflext.StructField) error {
	stmt, err := tb.dialect.CreateTable(tb.name, fields)
	if err != nil {
		return err
	}
	if _, err := sqldriver.Execute(
		tb.driver,
		stmt,
		tb.logger,
	); err != nil {
		return err
	}
	return nil
}

// we should skip virtual columns on insertion and migration
func skipGeneratedColumns(sfs []*reflext.StructField) (columns []string, fields []*reflext.StructField) {
	fields = make([]*reflext.StructField, 0, len(sfs))
	for _, sf := range sfs {
		if _, isOk := sf.Tag.LookUp("generated"); isOk {
			continue
		}
		columns = append(columns, sf.Path)
		fields = append(fields, sf)
	}
	return
}
