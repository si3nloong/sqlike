package logs

import (
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
)

// Logger :
type Logger interface {
	Debug(stmt *sqlstmt.Statement)
}
