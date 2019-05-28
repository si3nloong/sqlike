package mysql

import sqlstmt "github.com/si3nloong/sqlike/sqlike/sql/stmt"

// DropColumn :
func (ms *MySQL) DropColumn(table, column string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`ALTER TABLE ` + ms.Quote(table))
	stmt.WriteString(` DROP COLUMN ` + ms.Quote(column))
	stmt.WriteRune(';')
	return
}
