package common

import (
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
)

// Delete :
func (s *commonSQL) Delete(stmt db.Stmt, f *actions.DeleteActions) (err error) {
	err = buildStatement(stmt, s.parser, f)
	if err != nil {
		return
	}
	return
}
