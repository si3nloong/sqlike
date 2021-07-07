package debug

import (
	"github.com/si3nloong/sqlike/v2/sql/dialect"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
)

// ToSQL :
func ToSQL(src interface{}) error {
	ms := dialect.GetDialectByDriver("mysql")
	sqlstmt.NewStatement(ms)
	return nil
}
