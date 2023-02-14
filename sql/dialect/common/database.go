package common

import "github.com/si3nloong/sqlike/v2/db"

// UseDatabase :
func (s *commonSQL) UseDatabase(stmt db.Stmt, db string) {
	stmt.WriteString(`USE ` + s.Quote(db) + `;`)
}

// CreateDatabase :
func (s *commonSQL) CreateDatabase(stmt db.Stmt, db string, checkExists bool) {
	stmt.WriteString(`CREATE DATABASE `)
	if checkExists {
		stmt.WriteString(`IF NOT EXISTS `)
	}
	stmt.WriteString(s.Quote(db) + `;`)
}

// DropDatabase :
func (s *commonSQL) DropDatabase(stmt db.Stmt, db string, checkExists bool) {
	stmt.WriteString(`DROP SCHEMA `)
	if checkExists {
		stmt.WriteString(`IF EXISTS `)
	}
	stmt.WriteString(s.Quote(db) + `;`)
}

// GetDatabases :
func (s *commonSQL) GetDatabases(stmt db.Stmt) {
	stmt.WriteString(`SHOW DATABASES;`)
}
