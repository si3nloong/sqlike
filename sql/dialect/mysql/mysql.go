package mysql

import (
	"github.com/si3nloong/sqlike/db"
	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/sql/schema"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	sqlutil "github.com/si3nloong/sqlike/sql/util"
)

// MySQL :
type MySQL struct {
	schema *schema.Builder
	parser *sqlstmt.StatementBuilder
	sqlutil.MySQLUtil
}

var _ db.Dialect = (*(MySQL))(nil)

// New :
func New() *MySQL {
	sb := schema.NewBuilder()
	pr := sqlstmt.NewStatementBuilder()

	mySQLSchema{}.SetBuilders(sb)
	mySQLBuilder{}.SetRegistryAndBuilders(codec.DefaultRegistry, pr)

	return &MySQL{
		schema: sb,
		parser: pr,
	}
}

// GetVersion :
func (ms MySQL) GetVersion(stmt sqlstmt.Stmt) {
	stmt.WriteString("SELECT VERSION();")
}
