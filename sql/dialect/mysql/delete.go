package mysql

import (
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
)

// Delete :
func (ms *mySQL) Delete(stmt db.Stmt, f actions.DeleteActions) (err error) {
	err = buildStatement(stmt, ms.parser, f)
	return
}
