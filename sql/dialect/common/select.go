package common

import (
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/internal/primitive"
)

// Select :
func (s *commonSQL) Select(stmt db.Stmt, f *actions.FindActions, lck primitive.Lock) (err error) {
	err = s.parser.BuildStatement(stmt, f)
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
func (s *commonSQL) SelectStmt(stmt db.Stmt, query any) (err error) {
	err = s.parser.BuildStatement(stmt, query)
	stmt.WriteByte(';')
	return
}
