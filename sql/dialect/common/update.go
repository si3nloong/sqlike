package common

import (
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/db"
)

// Update :
func (s *commonSQL) Update(stmt db.Stmt, f *actions.UpdateActions) (err error) {
	err = buildStatement(stmt, s.parser, f)
	return
}
