package debug

import (
	"github.com/si3nloong/sqlike/v2/sql/dialect"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
)

// ToSQL :
func ToSQL(src any) error {
	ms := dialect.GetDialectByDriver("mysql")
	stmt := sqlstmt.AcquireStmt(ms)
	defer sqlstmt.ReleaseStmt(stmt)
	return nil
}
