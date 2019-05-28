package util

import "strconv"

// PostgresUtil :
type PostgresUtil struct{}

// Quote :
func (util PostgresUtil) Quote(n string) string {
	return strconv.Quote(n)
}
