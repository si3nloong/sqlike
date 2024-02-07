package mysql

import (
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql"
)

// Replace :
func (ms mySQL) Replace(stmt db.Stmt, db, table string, columns []string, query *sql.SelectStmt) (err error) {
	stmt.WriteString("REPLACE INTO ")
	stmt.WriteString(ms.TableName(db, table) + " ")
	if len(columns) > 0 {
		stmt.WriteByte('(')
		for i, col := range columns {
			if i > 0 {
				stmt.WriteByte(',')
			}
			stmt.WriteString(ms.Quote(col))
		}
		stmt.WriteByte(')')
		stmt.WriteByte(' ')
	}
	err = ms.parser.BuildStatement(stmt, query)
	if err != nil {
		return
	}
	stmt.WriteByte(';')
	return
}
