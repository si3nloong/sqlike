package sqldialect

import (
	"sync"

	"github.com/si3nloong/sqlike/reflext"
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
	GetVersion() (stmt *sqlstmt.Statement)
	GetDatabases() (stmt *sqlstmt.Statement)
	Format(v interface{}) (val string)
	HasTable(dbName, table string) (stmt *sqlstmt.Statement)
	RenameTable(oldName, newName string) (stmt *sqlstmt.Statement)
	DropColumn(table, column string) (stmt *sqlstmt.Statement)
	DropTable(table string, exists bool) (stmt *sqlstmt.Statement)
	TruncateTable(table string) (stmt *sqlstmt.Statement)
	GetColumns(dbName, table string) (stmt *sqlstmt.Statement)
	HasIndex(dbName, table, indexName string) (stmt *sqlstmt.Statement)
	GetIndexes(dbName, table string) (stmt *sqlstmt.Statement)
	CreateIndexes(table string, idxs []indexes.Index, supportDesc bool) (stmt *sqlstmt.Statement)
	DropIndex(table, idxName string) (stmt *sqlstmt.Statement)
	CreateTable(table, pk string, fields []*reflext.StructField) (stmt *sqlstmt.Statement, err error)
	AlterTable(table, pk string, fields []*reflext.StructField, columns util.StringSlice, indexes util.StringSlice, unsafe bool) (stmt *sqlstmt.Statement, err error)
	// ReplaceInto(table string, filter *sql.Query) (stmt string, args []interface{}, err error)
	InsertInto(table, pk string, columns []string, values [][]interface{}, opts *options.InsertOptions) (stmt *sqlstmt.Statement)
	Select(*actions.FindActions, options.LockMode) (stmt *sqlstmt.Statement, err error)
	Update(*actions.UpdateActions) (stmt *sqlstmt.Statement, err error)
	Delete(*actions.DeleteActions) (stmt *sqlstmt.Statement, err error)
	Copy(table string, columns []string, act *actions.CopyActions) (stmt *sqlstmt.Statement, err error)
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
