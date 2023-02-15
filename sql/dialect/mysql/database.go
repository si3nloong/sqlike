package mysql

import "github.com/si3nloong/sqlike/v2/db"

// UseDatabase :
func (s *mySQL) UseDatabase(stmt db.Stmt, db string) {
	stmt.WriteString(`USE ` + s.Quote(db) + `;`)
}

// CreateDatabase :
func (s *mySQL) CreateDatabase(stmt db.Stmt, db string, ifExists bool) {
	stmt.WriteString(`CREATE DATABASE `)
	if ifExists {
		stmt.WriteString(`IF NOT EXISTS `)
	}
	stmt.WriteString(s.Quote(db) + `;`)
}

// DropDatabase :
func (s *mySQL) DropDatabase(stmt db.Stmt, db string, ifExists bool) {
	stmt.WriteString(`DROP SCHEMA `)
	if ifExists {
		stmt.WriteString(`IF EXISTS `)
	}
	stmt.WriteString(s.Quote(db) + `;`)
}

// GetDatabases :
func (s *mySQL) GetDatabases(stmt db.Stmt) {
	stmt.WriteString(`SHOW DATABASES;`)
}
