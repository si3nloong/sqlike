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
	UseDatabase(db string) (stmt *sqlstmt.Statement)
	GetVersion() (stmt *sqlstmt.Statement)
	GetDatabases() (stmt *sqlstmt.Statement)
	CreateDatabase(db string, checkExists bool) (stmt *sqlstmt.Statement)
	DropDatabase(db string, checkExists bool) (stmt *sqlstmt.Statement)
	Format(v interface{}) (val string)
	HasTable(db, table string) (stmt *sqlstmt.Statement)
	HasPrimaryKey(db, table string) (stmt *sqlstmt.Statement)
	RenameTable(db, oldName, newName string) (stmt *sqlstmt.Statement)
	RenameColumn(db, table, oldColName, newColName string) (stmt *sqlstmt.Statement)
	DropColumn(db, table, column string) (stmt *sqlstmt.Statement)
	DropTable(db, table string, checkExists bool) (stmt *sqlstmt.Statement)
	TruncateTable(db, table string) (stmt *sqlstmt.Statement)
	GetColumns(db, table string) (stmt *sqlstmt.Statement)
	HasIndexByName(db, table, indexName string) (stmt *sqlstmt.Statement)
	HasIndex(dbName, table string, idx indexes.Index) (stmt *sqlstmt.Statement)
	GetIndexes(db, table string) (stmt *sqlstmt.Statement)
	CreateIndexes(db, table string, idxs []indexes.Index, supportDesc bool) (stmt *sqlstmt.Statement)
	DropIndexes(db, table string, idxs []string) (stmt *sqlstmt.Statement)
	CreateTable(db, table, pk string, info driver.Info, fields []*reflext.StructField) (stmt *sqlstmt.Statement, err error)
	AlterTable(db, table, pk string, hasPk bool, info driver.Info, fields []*reflext.StructField, columns util.StringSlice, indexes util.StringSlice, unsafe bool) (stmt *sqlstmt.Statement, err error)
	InsertInto(db, table, pk string, mapper *reflext.Mapper, codec codec.Codecer, fields []*reflext.StructField, values reflect.Value, opts *options.InsertOptions) (stmt *sqlstmt.Statement, err error)
	Select(*actions.FindActions, options.LockMode) (stmt *sqlstmt.Statement, err error)
	Update(*actions.UpdateActions) (stmt *sqlstmt.Statement, err error)
	Delete(*actions.DeleteActions) (stmt *sqlstmt.Statement, err error)
	SelectStmt(query interface{}) (stmt *sqlstmt.Statement, err error)
	Replace(db, table string, columns []string, query *sql.SelectStmt) (stmt *sqlstmt.Statement, err error)
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
