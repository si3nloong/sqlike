package mysql

import (
	"github.com/si3nloong/sqlike/actions"
	"github.com/si3nloong/sqlike/db"
)

// Update :
func (ms *MySQL) Update(stmt db.Stmt, f *actions.UpdateActions) (err error) {
	err = buildStatement(stmt, ms.parser, f)
	if err != nil {
		return
	}
	return
}
