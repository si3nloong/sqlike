package debug

import (
	"github.com/si3nloong/sqlike/db"
	"github.com/si3nloong/sqlike/sql/dialect/mysql"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
)

// ToSQL :
func ToSQL(src interface{}) error {
	ms := db.GetDialectByDriver("mysql").(*mysql.MySQL)
	sqlstmt.NewStatement(ms)
	return nil
}
