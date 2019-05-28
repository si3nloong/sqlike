package mysql

import (
	"bitbucket.org/SianLoong/sqlike/sqlike/sql/internal"
	sqlstmt "bitbucket.org/SianLoong/sqlike/sqlike/sql/stmt"
	sqlutil "bitbucket.org/SianLoong/sqlike/sqlike/sql/util"
)

// MySQL :
type MySQL struct {
	schema *internal.SchemaBuilder
	parser *sqlstmt.StatementParser
	sqlutil.MySQLUtil
}

// New :
func New() *MySQL {
	sb := internal.NewSchemaBuilder()
	pr := sqlstmt.NewStatementParser()

	mySQLSchema{}.SetBuilders(sb)
	mySQLParser{}.SetParsers(pr)

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
