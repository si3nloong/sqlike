package mysql

import "github.com/si3nloong/sqlike/v2/db"

// GetColumns :
func (ms *mySQL) GetColumns(stmt db.Stmt, dbName, table string) {
	stmt.AppendArgs(`SELECT ORDINAL_POSITION, COLUMN_NAME, COLUMN_TYPE, COLUMN_DEFAULT, IS_NULLABLE,
	DATA_TYPE, CHARACTER_SET_NAME, COLLATION_NAME, COLUMN_COMMENT, EXTRA FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION;`, dbName, table)
}

// RenameColumn :
func (ms *mySQL) RenameColumn(stmt db.Stmt, db, table, oldColName, newColName string) {
	stmt.WriteString("ALTER TABLE " + ms.TableName(db, table))
	stmt.WriteString(" RENAME COLUMN " + ms.Quote(oldColName) + " TO " + ms.Quote(newColName))
	stmt.WriteByte(';')
}

// DropColumn :
func (ms *mySQL) DropColumn(stmt db.Stmt, db, table, column string) {
	stmt.WriteString("ALTER TABLE " + ms.TableName(db, table))
	stmt.WriteString(" DROP COLUMN " + ms.Quote(column))
	stmt.WriteByte(';')
}
