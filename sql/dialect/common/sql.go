package common

import (
	"strconv"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql/schema"

	sqlstmt "github.com/si3nloong/sqlike/v2/sql/stmt"
)

type commonSQL struct {
	schema *schema.Builder
	parser *sqlstmt.StatementBuilder
	db.Codecer
}

var _ db.Dialect = (*(commonSQL))(nil)

// New :
func New() *commonSQL {
	sb := schema.NewBuilder()
	pr := sqlstmt.NewStatementBuilder()

	// codec := buildDefaultRegistry()
	// mySQLSchema{}.SetBuilders(sb)
	// mySQLBuilder{}.SetRegistryAndBuilders(codec, pr)

	return &commonSQL{
		schema: sb,
		parser: pr,
		// Codecer: codec,
	}
}

// GetVersion :
func (s *commonSQL) GetVersion(stmt db.Stmt) {
	stmt.WriteString("SELECT VERSION();")
}

func (s *commonSQL) TableName(db, table string) string {
	return "`" + db + "`.`" + table + "`"
}

// Var :
func (s *commonSQL) Var(i int) string {
	return "$" + strconv.Itoa(i)
}

// Quote :
func (s *commonSQL) Quote(n string) string {
	return "`" + n + "`"
}

// Wrap :
func (s *commonSQL) Wrap(n string) string {
	return "'" + n + "'"
}

// WrapOnlyValue :
func (s *commonSQL) WrapOnlyValue(n string) string {
	return n
}
