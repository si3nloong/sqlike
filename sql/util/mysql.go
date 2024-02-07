package util

import (
	"regexp"

	"github.com/si3nloong/sqlike/v2/internal/util"
)

var sqlFuncRegexp = regexp.MustCompile(`\([\w\_]+\)$`)

// MySQLUtil :
type MySQLUtil struct{}

// TableName :
func (u MySQLUtil) TableName(db, table string) string {
	return "`" + util.EscapeString(db, '`') + "`.`" + util.EscapeString(table, '`') + "`"
}

// Var :
func (u MySQLUtil) Var(i int) string {
	return "?"
}

// Quote :
func (u MySQLUtil) Quote(s string) string {
	return "`" + util.EscapeString(s, '`') + "`"
}

// Wrap :
func (u MySQLUtil) Wrap(s string) string {
	return "'" + util.EscapeString(s, '\'') + "'"
}

// WrapOnlyValue :
func (u MySQLUtil) WrapOnlyValue(n string) string {
	// eg. CURRENT_TIMESTAMP(6)
	if sqlFuncRegexp.MatchString(n) {
		return n
	}
	return u.Wrap(n)
}
