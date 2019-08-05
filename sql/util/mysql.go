package util

import "strings"

// MySQLUtil :
type MySQLUtil struct{}

// Quote :
func (util MySQLUtil) TableName(db, table string) string {
	return "`" + db + "`.`" + table + "`"
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
