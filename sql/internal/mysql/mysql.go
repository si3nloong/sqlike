package mysql

import (
	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/sql/schema"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	sqlutil "github.com/si3nloong/sqlike/sql/util"
)

// MySQL :
type MySQL struct {
	schema *schema.Builder
	parser *sqlstmt.StatementBuilder
	sqlutil.MySQLUtil
}

// New :
func New() *MySQL {
	sb := schema.NewBuilder()
	pr := sqlstmt.NewStatementBuilder()

	mySQLSchema{}.SetBuilders(sb)
	mySQLBuilder{}.SetRegistryAndBuilders(codec.DefaultRegistry, pr)

	return &MySQL{
		schema: sb,
		parser: pr,
	}
}

// GetVersion :
func (ms MySQL) GetVersion() (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`SELECT VERSION();`)
	return
}

// CreateDatabase :
func (ms MySQL) CreateDatabase(db string, exists bool) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("CREATE DATABASE")
	if !exists {
		stmt.WriteString(" IF NOT EXISTS")
	}
	stmt.WriteByte(' ')
	stmt.WriteString(ms.Quote(db) + ";")
	return
}

// DropDatabase :
func (ms MySQL) DropDatabase(db string, exists bool) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("DROP SCHEMA")
	if exists {
		stmt.WriteString(" IF EXISTS")
	}
	stmt.WriteByte(' ')
	stmt.WriteString(ms.Quote(db) + ";")
	return
}

// GetDatabases :
func (ms MySQL) GetDatabases() (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`SHOW DATABASES;`)
	return
}

// GetColumns :
func (ms MySQL) GetColumns(dbName, table string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`SELECT ORDINAL_POSITION, COLUMN_NAME, COLUMN_TYPE, COLUMN_DEFAULT, IS_NULLABLE,
	DATA_TYPE, CHARACTER_SET_NAME, COLLATION_NAME, EXTRA FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;`)
	stmt.AppendArgs([]interface{}{dbName, table})
	return
}
