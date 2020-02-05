package util

import "strings"

// MySQLUtil :
type MySQLUtil struct{}

// TableName :
func (util MySQLUtil) TableName(db, table string) string {
	return "`" + db + "`.`" + table + "`"
}

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
	// TODO: regex to check the string with () symbols
	if strings.Contains(n, "(") {
		return n
	}
	return util.Wrap(n)
}
