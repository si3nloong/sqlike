package mysql

import (
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
)

// Update :
func (ms *mySQL) Update(stmt db.Stmt, f actions.UpdateActions) (err error) {
	err = buildStatement(stmt, ms.parser, f)
	return
}
