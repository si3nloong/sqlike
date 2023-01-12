package sqlcommon

import (
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql/schema"

	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	sqlutil "github.com/si3nloong/sqlike/v2/sql/util"
)

type SQL struct {
	schema *schema.Builder
	parser *sqlstmt.StatementBuilder
	sqlutil.MySQLUtil
	db.Codecer
}
