package mysql

import (
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// Select :
func (ms *MySQL) Select(stmt sqlstmt.Stmt, f *actions.FindActions, lck options.LockMode) (err error) {
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
func (ms *MySQL) SelectStmt(stmt sqlstmt.Stmt, query interface{}) (err error) {
	err = ms.parser.BuildStatement(stmt, query)
	stmt.WriteByte(';')
	return
}

func buildStatement(stmt sqlstmt.Stmt, parser *sqlstmt.StatementBuilder, f interface{}) error {
	if err := parser.BuildStatement(stmt, f); err != nil {
		return err
	}
	stmt.WriteByte(';')
	return nil
}
