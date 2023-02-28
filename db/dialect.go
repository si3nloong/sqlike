package db

import (
	"reflect"

	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/internal/primitive"
	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/sql/util"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

type SqlDriver interface {
	Var(i int) string
}

// SQL dialect which implement interfaces
type Dialect interface {
	Codecer
	SqlDriver
	TableName(db, table string) string
	Quote(n string) string
	Format(v any) (val string)
	Connect(opt *options.ConnectOptions) (connStr string)
	UseDatabase(stmt Stmt, db string)
	GetVersion(stmt Stmt)
	GetDatabases(stmt Stmt)
	CreateDatabase(stmt Stmt, db string, exists bool)
	DropDatabase(stmt Stmt, db string, exists bool)
	HasTable(stmt Stmt, db, table string)
	HasPrimaryKey(stmt Stmt, db, table string)
	RenameTable(stmt Stmt, db, oldName, newName string)
	RenameColumn(stmt Stmt, db, table, oldColName, newColName string)
	DropColumn(stmt Stmt, db, table, column string)
	DropTable(stmt Stmt, db, table string, exists bool, unsafe bool)
	TruncateTable(stmt Stmt, db, table string)
	GetColumns(stmt Stmt, db, table string)
	HasIndexByName(stmt Stmt, db, table, indexName string)
	HasIndex(stmt Stmt, dbName, table string, idx sql.Index)
	GetIndexes(stmt Stmt, db, table string)
	CreateIndexes(stmt Stmt, db, table string, idxs []sql.Index, supportDesc bool)
	DropIndexes(stmt Stmt, db, table string, idxs []string)
	CreateTable(stmt Stmt, db, table, pk string, info Info, fields []reflext.FieldInfo) (err error)
	AlterTable(stmt Stmt, db, table, pk string, hasPk bool, info Info, fields []reflext.FieldInfo, columns util.StringSlice, indexes util.StringSlice, unsafe bool) (err error)
	InsertInto(stmt Stmt, db, table, pk string, mapper reflext.StructMapper, fields []reflext.FieldInfo, values reflect.Value, opts *options.InsertOptions) (err error)
	Select(stmt Stmt, act actions.FindActions, lock primitive.Lock) (err error)
	Update(stmt Stmt, act actions.UpdateActions) (err error)
	Delete(stmt Stmt, act actions.DeleteActions) (err error)
	SelectStmt(stmt Stmt, query any) (err error)
	Replace(stmt Stmt, db, table string, columns []string, query *sql.SelectStmt) (err error)
}
