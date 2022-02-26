package sqlike

import (
	"context"
	"reflect"
	"strings"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/x/reflext"

	"github.com/si3nloong/sqlike/v2/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/v2/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
)

// Table :
type Table struct {
	// current database name
	dbName string

	// table name
	name string

	// default primary key
	pk string

	client *Client

	// sql driver
	driver sqldriver.Driver

	// sql dialect
	dialect dialect.Dialect

	// logger
	logger db.Logger
}

// Rename : rename the current table name to new table name
func (tb *Table) Rename(ctx context.Context, name string) error {
	stmt := sqlstmt.AcquireStmt(tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	tb.dialect.RenameTable(stmt, tb.dbName, tb.name, name)
	_, err := sqldriver.Execute(
		ctx,
		getDriverFromContext(ctx, tb.driver),
		stmt,
		tb.logger,
	)
	return err
}

// Exists : this will return true when the table exists in the database
func (tb *Table) Exists(ctx context.Context) bool {
	var count int
	stmt := sqlstmt.AcquireStmt(tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	tb.dialect.HasTable(stmt, tb.dbName, tb.name)
	if err := sqldriver.QueryRowContext(
		ctx,
		getDriverFromContext(ctx, tb.driver),
		stmt,
		tb.logger,
	).Scan(&count); err != nil {
		panic(err)
	}
	return count > 0
}

// Columns :
func (tb *Table) Columns() *ColumnView {
	return &ColumnView{tb: tb}
}

// ListColumns : list all the column of the table.
func (tb *Table) ListColumns(ctx context.Context) ([]Column, error) {
	stmt := sqlstmt.AcquireStmt(tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	tb.dialect.GetColumns(stmt, tb.dbName, tb.name)
	rows, err := sqldriver.Query(
		ctx,
		getDriverFromContext(ctx, tb.driver),
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
			&col.Charset,
			&col.Collation,
			&col.Comment,
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

// ListIndexes : list all the index of the table.
func (tb *Table) ListIndexes(ctx context.Context) ([]Index, error) {
	stmt := sqlstmt.AcquireStmt(tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	tb.dialect.GetIndexes(stmt, tb.dbName, tb.name)
	rows, err := sqldriver.Query(
		ctx,
		getDriverFromContext(ctx, tb.driver),
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
			&idx.IsUnique,
		); err != nil {
			return nil, err
		}
		idx.IsUnique = !idx.IsUnique
		idxs = append(idxs, idx)
	}
	return idxs, nil
}

// MustMigrate : this will ensure the migrate is complete, otherwise it will panic
func (tb Table) MustMigrate(ctx context.Context, entity interface{}) {
	err := tb.Migrate(ctx, entity)
	if err != nil {
		panic(err)
	}
}

// Migrate : migrate will create a new table follows by the definition of struct tag, alter when the table already exists
func (tb *Table) Migrate(ctx context.Context, entity interface{}) error {
	return tb.migrateOne(ctx, tb.client.cache, entity, false)
}

// UnsafeMigrate : unsafe migration will delete non-exist index and columns, beware when you use this
func (tb *Table) UnsafeMigrate(ctx context.Context, entity interface{}) error {
	return tb.migrateOne(ctx, tb.client.cache, entity, true)
}

// MustUnsafeMigrate : this will panic if it get error on unsafe migrate
func (tb *Table) MustUnsafeMigrate(ctx context.Context, entity interface{}) {
	err := tb.migrateOne(ctx, tb.client.cache, entity, true)
	if err != nil {
		panic(err)
	}
}

// Truncate : delete all the table data.
func (tb *Table) Truncate(ctx context.Context) (err error) {
	stmt := sqlstmt.AcquireStmt(tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	tb.dialect.TruncateTable(stmt, tb.dbName, tb.name)
	_, err = sqldriver.Execute(
		ctx,
		getDriverFromContext(ctx, tb.driver),
		stmt,
		tb.logger,
	)
	return
}

// DropIfExists : will drop the table only if it exists.
func (tb *Table) DropIfExists(ctx context.Context) (err error) {
	stmt := sqlstmt.AcquireStmt(tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	tb.dialect.DropTable(stmt, tb.dbName, tb.name, true, false)
	_, err = sqldriver.Execute(
		ctx,
		getDriverFromContext(ctx, tb.driver),
		stmt,
		tb.logger,
	)
	return
}

// Drop : drop the table, but it might throw error when the table is not exists
func (tb *Table) Drop(ctx context.Context) (err error) {
	stmt := sqlstmt.AcquireStmt(tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	tb.dialect.DropTable(stmt, tb.dbName, tb.name, false, false)
	_, err = sqldriver.Execute(
		ctx,
		getDriverFromContext(ctx, tb.driver),
		stmt,
		tb.logger,
	)
	return
}

// UnsafeDrop : drop the table without table is exists and foreign key constraint error
func (tb *Table) UnsafeDrop(ctx context.Context) (err error) {
	stmt := sqlstmt.AcquireStmt(tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	tb.dialect.DropTable(stmt, tb.dbName, tb.name, true, true)
	_, err = sqldriver.Execute(
		ctx,
		getDriverFromContext(ctx, tb.driver),
		stmt,
		tb.logger,
	)
	return
}

// Replace :
func (tb *Table) Replace(ctx context.Context, fields []string, query *sql.SelectStmt) error {
	stmt := sqlstmt.AcquireStmt(tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	if err := tb.dialect.Replace(
		stmt,
		tb.dbName,
		tb.name,
		fields,
		query,
	); err != nil {
		return err
	}

	if _, err := sqldriver.Execute(
		ctx,
		getDriverFromContext(ctx, tb.driver),
		stmt,
		tb.logger,
	); err != nil {
		return err
	}
	return nil
}

// Indexes :
func (tb *Table) Indexes() *IndexView {
	return &IndexView{tb: tb}
}

// HasIndexByName :
func (tb *Table) HasIndexByName(ctx context.Context, name string) (bool, error) {
	return isIndexExists(
		ctx,
		tb.dbName,
		tb.name,
		name,
		tb.driver,
		tb.dialect,
		tb.logger,
	)
}

func (tb *Table) migrateOne(ctx context.Context, cache reflext.StructMapper, entity interface{}, unsafe bool) error {
	v := reflext.ValueOf(entity)
	if !v.IsValid() {
		return ErrInvalidInput
	}

	t := reflext.Deref(v.Type())
	if !reflext.IsKind(t, reflect.Struct) {
		return ErrExpectedStruct
	}

	cdc := cache.CodecByType(t)
	fields := skipColumns(cdc.Properties(), nil)
	if len(fields) < 1 {
		return ErrEmptyFields
	}

	if !tb.Exists(ctx) {
		return tb.createTable(ctx, fields)
	}

	columns, err := tb.ListColumns(ctx)
	if err != nil {
		return err
	}
	idxs, err := tb.ListIndexes(ctx)
	if err != nil {
		return err
	}
	return tb.alterTable(ctx, fields, columns, idxs, unsafe)
}

func (tb *Table) createTable(ctx context.Context, fields []reflext.FieldInfo) error {
	stmt := sqlstmt.AcquireStmt(tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	if err := tb.dialect.CreateTable(
		stmt,
		tb.dbName,
		tb.name,
		tb.pk,
		tb.client.DriverInfo,
		fields,
	); err != nil {
		return err
	}
	if _, err := sqldriver.Execute(
		ctx,
		getDriverFromContext(ctx, tb.driver),
		stmt,
		tb.logger,
	); err != nil {
		return err
	}
	return nil
}

func (tb *Table) alterTable(ctx context.Context, fields []reflext.FieldInfo, columns []Column, indexs []Index, unsafe bool) error {
	cols := make([]string, len(columns))
	for i, col := range columns {
		cols[i] = col.Name
	}
	idxs := make([]string, len(indexs))
	for i, idx := range indexs {
		idxs[i] = idx.Name
	}
	stmt := sqlstmt.AcquireStmt(tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	tb.dialect.HasPrimaryKey(stmt, tb.dbName, tb.name)
	var count uint
	if err := sqldriver.QueryRowContext(
		ctx,
		tb.driver,
		stmt,
		tb.logger,
	).Scan(&count); err != nil {
		return err
	}
	stmt.Reset()
	if err := tb.dialect.AlterTable(
		stmt,
		tb.dbName, tb.name, tb.pk, count > 0,
		tb.client.DriverInfo,
		fields, cols, idxs, unsafe,
	); err != nil {
		return err
	}
	if _, err := sqldriver.Execute(
		ctx,
		getDriverFromContext(ctx, tb.driver),
		stmt,
		tb.logger,
	); err != nil {
		return err
	}
	return nil
}
