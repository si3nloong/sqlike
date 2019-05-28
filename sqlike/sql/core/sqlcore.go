package sqlcore

import (
	"sync"

	"bitbucket.org/SianLoong/sqlike/reflext"
	"bitbucket.org/SianLoong/sqlike/sqlike/actions"
	"bitbucket.org/SianLoong/sqlike/sqlike/indexes"
	"bitbucket.org/SianLoong/sqlike/sqlike/options"
	"bitbucket.org/SianLoong/sqlike/sqlike/sql/core/mysql"
	sqlstmt "bitbucket.org/SianLoong/sqlike/sqlike/sql/stmt"
)

// Dialect :
type Dialect interface {
	GetVersion() (stmt *sqlstmt.Statement)
	GetDatabases() (stmt *sqlstmt.Statement)
	Format(v interface{}) (val string)
	DropColumn(table, column string) (stmt *sqlstmt.Statement)
	DropTable(table string, exists bool) (stmt *sqlstmt.Statement)
	TruncateTable(table string) (stmt *sqlstmt.Statement)
	HasTable(dbName, table string) (stmt *sqlstmt.Statement)
	GetColumns(dbName, table string) (stmt *sqlstmt.Statement)
	GetIndexes(dbName, table string) (stmt *sqlstmt.Statement)
	CreateIndexes(table string, idxs []indexes.Index) (stmt *sqlstmt.Statement)
	DropIndex(table, idxName string) (stmt *sqlstmt.Statement)
	CreateTable(table string, fields []*reflext.StructField) (stmt *sqlstmt.Statement, err error)
	AlterTable(table string, fields []*reflext.StructField, columns []string, indexes []string, unsafe bool) (stmt *sqlstmt.Statement, err error)
	// ReplaceInto(table string, filter *sql.Query) (stmt string, args []interface{}, err error)
	InsertInto(table string, columns []string, values [][]interface{}, opts *options.InsertOptions) (stmt *sqlstmt.Statement)
	Select(*actions.FindActions) (stmt *sqlstmt.Statement, err error)
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
