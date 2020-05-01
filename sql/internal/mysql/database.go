package mysql

import (
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
)

// UseDatabase :
func (ms MySQL) UseDatabase(stmt sqlstmt.Stmt, db string) {
	stmt.WriteString("USE " + ms.Quote(db) + ";")
}

// CreateDatabase :
func (ms MySQL) CreateDatabase(stmt sqlstmt.Stmt, db string, checkExists bool) {
	stmt.WriteString("CREATE DATABASE")
	if checkExists {
		stmt.WriteString(" IF NOT EXISTS")
	}
	stmt.WriteByte(' ')
	stmt.WriteString(ms.Quote(db) + ";")
}

// DropDatabase :
func (ms MySQL) DropDatabase(stmt sqlstmt.Stmt, db string, checkExists bool) {
	stmt.WriteString("DROP SCHEMA")
	if checkExists {
		stmt.WriteString(" IF EXISTS")
	}
	stmt.WriteByte(' ')
	stmt.WriteString(ms.Quote(db) + ";")
}

// GetDatabases :
func (ms MySQL) GetDatabases(stmt sqlstmt.Stmt) {
	stmt.WriteString("SHOW DATABASES;")
}
