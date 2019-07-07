package mysql

import (
	"strings"

	"github.com/si3nloong/sqlike/sqlike/options"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
)

// InsertInto :
func (ms MySQL) InsertInto(table, pk string, columns []string, values [][]interface{}, opt *options.InsertOptions) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`INSERT`)
	if opt.Mode == options.InsertIgnore {
		stmt.WriteString(` IGNORE`)
	}
	stmt.WriteString(` INTO ` + ms.Quote(table) + ` (`)
	for i, col := range columns {
		if i > 0 {
			stmt.WriteRune(',')
		}
		stmt.WriteString(ms.Quote(col))
	}
	stmt.WriteString(`) VALUES `)
	binds := strings.Repeat(`?,`, len(values[0]))
	binds = binds[:len(binds)-1]
	length := len(values)
	for i := 0; i < length; i++ {
		if i > 0 {
			stmt.WriteRune(',')
		}
		stmt.WriteString("(" + binds[:] + ")")
		stmt.AppendArgs(values[0])
		values = values[1:]
	}
	if opt.Mode == options.InsertOnDuplicate {
		stmt.WriteString(` ON DUPLICATE KEY UPDATE `)
		next := false
		for _, col := range columns {
			if col == pk {
				next = false
				continue
			}
			if next {
				stmt.WriteRune(',')
			}
			c := ms.Quote(col)
			stmt.WriteString(c + `=VALUES(` + c + `)`)
			next = true
		}
	}
	stmt.WriteRune(';')
	return
}
