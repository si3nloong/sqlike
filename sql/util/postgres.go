package util

import "strconv"

// PostgresUtil :
type PostgresUtil struct{}

// TableName :
func (util PostgresUtil) TableName(db, table string) string {
	return "`" + db + "`.`" + table + "`"
}

// Var :
func (util PostgresUtil) Var(i int) string {
	return "$" + strconv.Itoa(i)
}

// Quote :
func (util PostgresUtil) Quote(n string) string {
	return strconv.Quote(n)
}

// Wrap :
func (util PostgresUtil) Wrap(n string) string {
	return "'" + n + "'"
}
