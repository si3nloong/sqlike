package mysql

import (
	"github.com/si3nloong/sqlike/actions"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
)

// Update :
func (ms *MySQL) Update(stmt sqlstmt.Stmt, f *actions.UpdateActions) (err error) {
	err = buildStatement(stmt, ms.parser, f)
	if err != nil {
		return
	}
	return
}
