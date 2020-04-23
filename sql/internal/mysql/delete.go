package mysql

import (
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/actions"
)

// Delete :
func (ms *MySQL) Delete(f *actions.DeleteActions) (stmt *sqlstmt.Statement, err error) {
	stmt = sqlstmt.NewStatement(ms)
	err = buildStatement(stmt, ms.parser, f)
	if err != nil {
		return
	}
	return
}
