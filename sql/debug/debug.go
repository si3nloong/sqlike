package debug

import (
	"github.com/si3nloong/sqlike/v2/sql/dialect"
	"github.com/si3nloong/sqlike/v2/sql/dialect/mysql"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
)

// ToSQL :
func ToSQL(src interface{}) error {
	ms := dialect.GetDialectByDriver("mysql").(*mysql.MySQL)
	sqlstmt.NewStatement(ms)
	return nil
}
