package util

import "regexp"

var sqlFuncRegexp = regexp.MustCompile(`\([\w\_]+\)$`)

// MySQLUtil :
type MySQLUtil struct{}

// TableName :
func (util MySQLUtil) TableName(db, table string) string {
	return "`" + db + "`.`" + table + "`"
}

// Var :
func (util MySQLUtil) Var(i int) string {
	return "?"
}

// Quote :
func (util MySQLUtil) Quote(n string) string {
	return "`" + n + "`"
}

// Wrap :
func (util MySQLUtil) Wrap(n string) string {
	return "'" + n + "'"
}

// WrapOnlyValue :
func (util MySQLUtil) WrapOnlyValue(n string) string {
	// eg. CURRENT_TIMESTAMP(6)
	if sqlFuncRegexp.MatchString(n) {
		return n
	}
	return util.Wrap(n)
}
