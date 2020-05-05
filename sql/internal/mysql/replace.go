package mysql

import (
	"github.com/si3nloong/sqlike/sql"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
)

// Replace :
func (ms MySQL) Replace(stmt sqlstmt.Stmt, db, table string, columns []string, query *sql.SelectStmt) (err error) {
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
