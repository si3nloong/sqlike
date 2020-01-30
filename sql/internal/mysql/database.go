package mysql

import sqlstmt "github.com/si3nloong/sqlike/sql/stmt"

// UseDatabase :
func (ms MySQL) UseDatabase(db string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("USE " + ms.Quote(db) + ";")
	return
}

// CreateDatabase :
func (ms MySQL) CreateDatabase(db string, checkExists bool) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("CREATE DATABASE")
	if checkExists {
		stmt.WriteString(" IF NOT EXISTS")
	}
	stmt.WriteByte(' ')
	stmt.WriteString(ms.Quote(db) + ";")
	return
}

// DropDatabase :
func (ms MySQL) DropDatabase(db string, checkExists bool) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("DROP SCHEMA")
	if checkExists {
		stmt.WriteString(" IF EXISTS")
	}
	stmt.WriteByte(' ')
	stmt.WriteString(ms.Quote(db) + ";")
	return
}

// GetDatabases :
func (ms MySQL) GetDatabases() (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("SHOW DATABASES;")
	return
}
