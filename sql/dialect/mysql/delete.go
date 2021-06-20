package mysql

import (
	"github.com/si3nloong/sqlike/actions"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
)

// Delete :
func (ms *MySQL) Delete(stmt sqlstmt.Stmt, f *actions.DeleteActions) (err error) {
	err = buildStatement(stmt, ms.parser, f)
	if err != nil {
		return
	}
	return
}
