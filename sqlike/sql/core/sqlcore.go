package sqlcore

import (
	"sync"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/sqlike/sql/core/mysql"
	sqlstmt "github.com/si3nloong/sqlike/sqlike/sql/stmt"
)

// Dialect :
type Dialect interface {
	Connect(opt *options.ConnectOptions) (connStr string)
	GetVersion() (stmt *sqlstmt.Statement)
	GetDatabases() (stmt *sqlstmt.Statement)
	Format(v interface{}) (val string)
	HasTable(dbName, table string) (stmt *sqlstmt.Statement)
	RenameTable(oldName, newName string) (stmt *sqlstmt.Statement)
	DropColumn(table, column string) (stmt *sqlstmt.Statement)
	DropTable(table string, exists bool) (stmt *sqlstmt.Statement)
	TruncateTable(table string) (stmt *sqlstmt.Statement)
	GetColumns(dbName, table string) (stmt *sqlstmt.Statement)
	GetIndexes(dbName, table string) (stmt *sqlstmt.Statement)
	CreateIndexes(table string, idxs []indexes.Index) (stmt *sqlstmt.Statement)
	DropIndex(table, idxName string) (stmt *sqlstmt.Statement)
	CreateTable(table string, fields []*reflext.StructField) (stmt *sqlstmt.Statement, err error)
	AlterTable(table string, fields []*reflext.StructField, columns []string, indexes []string, unsafe bool) (stmt *sqlstmt.Statement, err error)
	// ReplaceInto(table string, filter *sql.Query) (stmt string, args []interface{}, err error)
	InsertInto(table string, columns []string, values [][]interface{}, opts *options.InsertOptions) (stmt *sqlstmt.Statement)
	Select(*actions.FindActions, options.LockMode) (stmt *sqlstmt.Statement, err error)
	Update(*actions.UpdateActions) (stmt *sqlstmt.Statement, err error)
	Delete(*actions.DeleteActions) (stmt *sqlstmt.Statement, err error)
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
