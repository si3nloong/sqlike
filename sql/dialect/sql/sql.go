package sql

import (
	"strconv"

	"github.com/si3nloong/sqlike/v2/internal/util"
)

type CommonSql struct{}

// ## MySQL
// Quotes (Single and Double) are used around strings. However, backticks are used around table and column identifiers.
// https://chartio.com/learn/sql-tips/single-double-quote-and-backticks-in-mysql-queries/
//
// ## Postgres
// Double quotes are used to indicate identifiers within the database, which are objects like tables, column names, and roles. In contrast, single quotes are used to indicate string literals.

// TableName :
func (s CommonSql) TableName(db, table string) string {
	return strconv.Quote(db) + "." + strconv.Quote(table)
}

// Var :
func (s CommonSql) Var(i int) string {
	return "$" + strconv.Itoa(i)
}

// Quote :
func (s CommonSql) Quote(v string) string {
	return "`" + util.EscapeString(v, '`') + "`"
}

// Wrap :
func (s CommonSql) Wrap(v string) string {
	return "'" + util.EscapeString(v, '\'') + "'"
}
