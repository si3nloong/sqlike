package mysql

import (
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql/schema"
	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
	sqlutil "github.com/si3nloong/sqlike/v2/sql/util"
)

type mySQL struct {
	schema *schema.Builder
	parser *sqlstmt.StatementBuilder
	sqlutil.MySQLUtil
	db.Codecer
}

var _ db.Dialect = (*(mySQL))(nil)

// New :
func New() *mySQL {
	sb := schema.NewBuilder()
	pr := sqlstmt.NewStatementBuilder()

	codec := buildDefaultRegistry()
	mySQLSchema{}.SetBuilders(sb)
	mySQLBuilder{}.SetRegistryAndBuilders(codec, pr)

	return &mySQL{
		schema:  sb,
		parser:  pr,
		Codecer: codec,
	}
}

// GetVersion :
func (ms mySQL) GetVersion(stmt db.Stmt) {
	stmt.WriteString("SELECT VERSION();")
}
