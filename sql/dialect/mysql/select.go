package mysql

import (
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/options"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
)

// Select :
func (ms *MySQL) Select(stmt db.Stmt, f *actions.FindActions, lck options.LockMode) (err error) {
	err = ms.parser.BuildStatement(stmt, f)
	if err != nil {
		return
	}
	switch lck {
	case options.LockForUpdate:
		stmt.WriteString(" FOR UPDATE")
	case options.LockForRead:
		stmt.WriteString(" LOCK IN SHARE MODE")
	}
	stmt.WriteByte(';')
	return
}

// SelectStmt :
func (ms *MySQL) SelectStmt(stmt db.Stmt, query interface{}) (err error) {
	err = ms.parser.BuildStatement(stmt, query)
	stmt.WriteByte(';')
	return
}

func buildStatement(stmt db.Stmt, parser *sqlstmt.StatementBuilder, f interface{}) error {
	if err := parser.BuildStatement(stmt, f); err != nil {
		return err
	}
	stmt.WriteByte(';')
	return nil
}
