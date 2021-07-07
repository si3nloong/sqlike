package mysql

import "github.com/si3nloong/sqlike/v2/db"

// UseDatabase :
func (ms mySQL) UseDatabase(stmt db.Stmt, db string) {
	stmt.WriteString("USE " + ms.Quote(db) + ";")
}

// CreateDatabase :
func (ms mySQL) CreateDatabase(stmt db.Stmt, db string, checkExists bool) {
	stmt.WriteString("CREATE DATABASE")
	if checkExists {
		stmt.WriteString(" IF NOT EXISTS")
	}
	stmt.WriteByte(' ')
	stmt.WriteString(ms.Quote(db) + ";")
}

// DropDatabase :
func (ms mySQL) DropDatabase(stmt db.Stmt, db string, checkExists bool) {
	stmt.WriteString("DROP SCHEMA")
	if checkExists {
		stmt.WriteString(" IF EXISTS")
	}
	stmt.WriteByte(' ')
	stmt.WriteString(ms.Quote(db) + ";")
}

// GetDatabases :
func (ms mySQL) GetDatabases(stmt db.Stmt) {
	stmt.WriteString("SHOW DATABASES;")
}
