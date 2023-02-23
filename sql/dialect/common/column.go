package common

import "github.com/si3nloong/sqlike/v2/db"

// GetColumns :
func (s *commonSQL) GetColumns(stmt db.Stmt, dbName, table string) {
	stmt.AppendArgs(`SELECT ORDINAL_POSITION, COLUMN_NAME, COLUMN_TYPE, COLUMN_DEFAULT, IS_NULLABLE, DATA_TYPE, CHARACTER_SET_NAME, COLLATION_NAME, COLUMN_COMMENT, EXTRA FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = `+s.Var(stmt.Pos()+1)+` AND TABLE_NAME = `+s.Var(stmt.Pos()+1)+` ORDER BY ORDINAL_POSITION;`, dbName, table)
}

// RenameColumn :
func (s *commonSQL) RenameColumn(stmt db.Stmt, db, table, oldColName, newColName string) {
	stmt.WriteString(`ALTER TABLE ` + s.TableName(db, table) + ` RENAME COLUMN ` + s.Quote(oldColName) + ` TO ` + s.Quote(newColName) + `;`)
}

// DropColumn :
func (s *commonSQL) DropColumn(stmt db.Stmt, db, table, column string) {
	stmt.WriteString(`ALTER TABLE ` + s.TableName(db, table) + ` DROP COLUMN ` + s.Quote(column) + `;`)
}
