package mysql

import (
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// Select :
func (ms *MySQL) Select(f *actions.FindActions, lck options.LockMode) (stmt *sqlstmt.Statement, err error) {
	stmt = sqlstmt.NewStatement(ms)
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
	stmt.WriteRune(';')
	return
}

// SelectStmt :
func (ms *MySQL) SelectStmt(query interface{}) (stmt *sqlstmt.Statement, err error) {
	stmt = sqlstmt.NewStatement(ms)
	err = ms.parser.BuildStatement(stmt, query)
	stmt.WriteRune(';')
	return
}

func buildStatement(stmt *sqlstmt.Statement, parser *sqlstmt.StatementBuilder, f interface{}) error {
	if err := parser.BuildStatement(stmt, f); err != nil {
		return err
	}
	stmt.WriteRune(';')
	return nil
}
