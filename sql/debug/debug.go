package debug

import (
	"log"

	"github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sql/internal/mysql"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
)

func init() {
	ToSQL(nil)
}

// ToSQL :
func ToSQL(src interface{}) error {
	ms := dialect.GetDialectByDriver("mysql").(*mysql.MySQL)
	stmt := sqlstmt.NewStatement(ms)
	log.Println(stmt)
	return nil
}
