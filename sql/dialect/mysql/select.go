package mysql

import (
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	"github.com/si3nloong/sqlike/v2/x/primitive"
)

// Select :
func (ms *mySQL) Select(stmt db.Stmt, f *actions.FindActions, lck primitive.Lock) (err error) {
	err = ms.parser.BuildStatement(stmt, f)
	if err != nil {
		return
	}
	switch lck.Type {
	case primitive.LockForUpdate:
		stmt.WriteString(" FOR UPDATE")
	case primitive.LockForShare:
		stmt.WriteString(" LOCK IN SHARE MODE")
	}
	stmt.WriteByte(';')
	return
}

// SelectStmt :
func (ms *mySQL) SelectStmt(stmt db.Stmt, query interface{}) (err error) {
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
