package postgres

import (
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql/schema"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	sqlutil "github.com/si3nloong/sqlike/v2/sql/util"
)

type postgreSQL struct {
	schema *schema.Builder
	parser *sqlstmt.StatementBuilder
	sqlutil.PostgresUtil
	db.Codecer
}
