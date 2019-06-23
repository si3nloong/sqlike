package logs

import (
	sqlstmt "github.com/si3nloong/sqlike/sqlike/sql/stmt"
)

// Logger :
type Logger interface {
	Format(stmt *sqlstmt.Statement)
}
