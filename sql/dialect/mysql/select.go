package mysql

import (
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/internal/primitive"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
)

// Select :
func (ms *mySQL) Select(stmt db.Stmt, f actions.FindActions, lck primitive.Lock) (err error) {
	err = ms.parser.BuildStatement(stmt, f)
	if err != nil {
		return
	}
	// if lck != nil {
	switch lck.Type {
	case primitive.LockForUpdate:
		stmt.WriteString(` FOR UPDATE`)
	case primitive.LockForShare:
		stmt.WriteString(` FOR SHARE`)
	}
	if lck.Of != nil {
		ms.parser.BuildStatement(stmt, *lck.Of)
	}
	switch lck.Option {
	case primitive.NoWait:
		stmt.WriteString(` NOWAIT`)
	case primitive.SkipLocked:
		stmt.WriteString(` SKIP LOCKED`)
	}
	// }
	stmt.WriteByte(';')
	return
}

// SelectStmt :
func (ms *mySQL) SelectStmt(stmt db.Stmt, query any) (err error) {
	err = ms.parser.BuildStatement(stmt, query)
	stmt.WriteByte(';')
	return
}

func buildStatement(stmt db.Stmt, parser *sqlstmt.StatementBuilder, f any) error {
	if err := parser.BuildStatement(stmt, f); err != nil {
		return err
	}
	stmt.WriteByte(';')
	return nil
}
