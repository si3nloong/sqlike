package mysql

import (
	"github.com/si3nloong/sqlike/actions"
	"github.com/si3nloong/sqlike/db"
)

// Delete :
func (ms *MySQL) Delete(stmt db.Stmt, f *actions.DeleteActions) (err error) {
	err = buildStatement(stmt, ms.parser, f)
	if err != nil {
		return
	}
	return
}
