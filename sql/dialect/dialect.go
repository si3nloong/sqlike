package dialect

import (
	"reflect"
	"strings"
	"sync"

	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/sql/driver"
	"github.com/si3nloong/sqlike/v2/sql/util"
	"github.com/si3nloong/sqlike/v2/x/primitive"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// Dialect :
type Dialect interface {
	db.Codecer
	TableName(db, table string) string
	Var(i int) string
	Quote(n string) string
	Format(v interface{}) (val string)
	Connect(opt *options.ConnectOptions) (connStr string)
	UseDatabase(stmt db.Stmt, db string)
	GetVersion(stmt db.Stmt)
	GetDatabases(stmt db.Stmt)
	CreateDatabase(stmt db.Stmt, db string, checkExists bool)
	DropDatabase(stmt db.Stmt, db string, checkExists bool)
	HasTable(stmt db.Stmt, db, table string)
	HasPrimaryKey(stmt db.Stmt, db, table string)
	RenameTable(stmt db.Stmt, db, oldName, newName string)
	RenameColumn(stmt db.Stmt, db, table, oldColName, newColName string)
	DropColumn(stmt db.Stmt, db, table, column string)
	DropTable(stmt db.Stmt, db, table string, checkExists bool, unsafe bool)
	TruncateTable(stmt db.Stmt, db, table string)
	GetColumns(stmt db.Stmt, db, table string)
	HasIndexByName(stmt db.Stmt, db, table, indexName string)
	HasIndex(stmt db.Stmt, dbName, table string, idx sql.Index)
	GetIndexes(stmt db.Stmt, db, table string)
	CreateIndexes(stmt db.Stmt, db, table string, idxs []sql.Index, supportDesc bool)
	DropIndexes(stmt db.Stmt, db, table string, idxs []string)
	CreateTable(stmt db.Stmt, db, table, pk string, info driver.Info, fields []reflext.FieldInfo) (err error)
	AlterTable(stmt db.Stmt, db, table, pk string, hasPk bool, info driver.Info, fields []reflext.FieldInfo, columns util.StringSlice, indexes util.StringSlice, unsafe bool) (err error)
	InsertInto(stmt db.Stmt, db, table, pk string, mapper reflext.StructMapper, fields []reflext.FieldInfo, values reflect.Value, opts *options.InsertOptions) (err error)
	Select(stmt db.Stmt, act *actions.FindActions, lock primitive.Lock) (err error)
	Update(stmt db.Stmt, act *actions.UpdateActions) (err error)
	Delete(stmt db.Stmt, act *actions.DeleteActions) (err error)
	SelectStmt(stmt db.Stmt, query interface{}) (err error)
	Replace(stmt db.Stmt, db, table string, columns []string, query *sql.SelectStmt) (err error)
}

var (
	mutex    sync.Mutex
	dialects = make(map[string]Dialect)
)

// RegisterDialect :
func RegisterDialect(driver string, dialect Dialect) {
	mutex.Lock()
	defer mutex.Unlock()
	if dialect == nil {
		panic("invalid nil dialect")
	}
	dialects[driver] = dialect
}

// GetDialectByDriver :
func GetDialectByDriver(driver string) Dialect {
	driver = strings.TrimSpace(strings.ToLower(driver))
	return dialects[driver]
}
