package mysql

import sqlstmt "github.com/si3nloong/sqlike/sql/stmt"

// DropColumn :
func (ms *MySQL) DropColumn(db, table, column string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("ALTER TABLE " + ms.TableName(db, table))
	stmt.WriteString(" DROP COLUMN " + ms.Quote(column))
	stmt.WriteRune(';')
	return
}
