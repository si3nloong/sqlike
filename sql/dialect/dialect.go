package dialect

import (
	"reflect"
	"sync"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql"
	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sql/internal/mysql"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sql/util"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// Dialect :
type Dialect interface {
	Connect(opt *options.ConnectOptions) (connStr string)
	UseDatabase(stmt sqlstmt.Stmt, db string)
	GetVersion(stmt sqlstmt.Stmt)
	GetDatabases(stmt sqlstmt.Stmt)
	CreateDatabase(stmt sqlstmt.Stmt, db string, checkExists bool)
	DropDatabase(stmt sqlstmt.Stmt, db string, checkExists bool)
	HasTable(stmt sqlstmt.Stmt, db, table string)
	HasPrimaryKey(stmt sqlstmt.Stmt, db, table string)
	RenameTable(stmt sqlstmt.Stmt, db, oldName, newName string)
	RenameColumn(stmt sqlstmt.Stmt, db, table, oldColName, newColName string)
	DropColumn(stmt sqlstmt.Stmt, db, table, column string)
	DropTable(stmt sqlstmt.Stmt, db, table string, checkExists bool)
	TruncateTable(stmt sqlstmt.Stmt, db, table string)
	GetColumns(stmt sqlstmt.Stmt, db, table string)
	HasIndexByName(stmt sqlstmt.Stmt, db, table, indexName string)
	HasIndex(stmt sqlstmt.Stmt, dbName, table string, idx indexes.Index)
	GetIndexes(stmt sqlstmt.Stmt, db, table string)
	CreateIndexes(stmt sqlstmt.Stmt, db, table string, idxs []indexes.Index, supportDesc bool)
	DropIndexes(stmt sqlstmt.Stmt, db, table string, idxs []string)
	CreateTable(stmt sqlstmt.Stmt, db, table, pk string, info driver.Info, fields []*reflext.StructField) (err error)
	AlterTable(stmt sqlstmt.Stmt, db, table, pk string, hasPk bool, info driver.Info, fields []*reflext.StructField, columns util.StringSlice, indexes util.StringSlice, unsafe bool) (err error)
	InsertInto(stmt sqlstmt.Stmt, db, table, pk string, mapper *reflext.Mapper, codec codec.Codecer, fields []*reflext.StructField, values reflect.Value, opts *options.InsertOptions) (err error)
	Select(*actions.FindActions, options.LockMode) (stmt *sqlstmt.Statement, err error)
	Update(*actions.UpdateActions) (stmt *sqlstmt.Statement, err error)
	Delete(*actions.DeleteActions) (stmt *sqlstmt.Statement, err error)
	SelectStmt(query interface{}) (stmt *sqlstmt.Statement, err error)
	Replace(db, table string, columns []string, query *sql.SelectStmt) (stmt *sqlstmt.Statement, err error)
	Var(i int) string
	Format(v interface{}) (val string)
}

var (
	mutex    sync.Mutex
	dialects = make(map[string]Dialect)
)

var _ Dialect = (*(mysql.MySQL))(nil)

func init() {
	RegisterDialect("mysql", mysql.New())
}

// RegisterDialect :
func RegisterDialect(driver string, dialect Dialect) {
	mutex.Lock()
	defer mutex.Unlock()
	dialects[driver] = dialect
}

// GetDialectByDriver :
func GetDialectByDriver(driver string) Dialect {
	return dialects[driver]
}
