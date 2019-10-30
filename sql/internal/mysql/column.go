package mysql

import sqlstmt "github.com/si3nloong/sqlike/sql/stmt"

// GetColumns :
func (ms *MySQL) GetColumns(dbName, table string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`SELECT ORDINAL_POSITION, COLUMN_NAME, COLUMN_TYPE, COLUMN_DEFAULT, IS_NULLABLE,
	DATA_TYPE, CHARACTER_SET_NAME, COLLATION_NAME, EXTRA FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;`)
	stmt.AppendArgs([]interface{}{dbName, table})
	return
}

// DropColumn :
func (ms *MySQL) DropColumn(db, table, column string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("ALTER TABLE " + ms.TableName(db, table))
	stmt.WriteString(" DROP COLUMN " + ms.Quote(column))
	stmt.WriteRune(';')
	return
}
