package common

import (
	"github.com/si3nloong/sqlike/v2/db"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
)

func buildStatement(stmt db.Stmt, parser *sqlstmt.StatementBuilder, f any) error {
	if err := parser.BuildStatement(stmt, f); err != nil {
		return err
	}
	stmt.WriteByte(';')
	return nil
}
