package common

import (
	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql"
)

// Replace :
func (s *commonSQL) Replace(stmt db.Stmt, db, table string, columns []string, query *sql.SelectStmt) (err error) {
	stmt.WriteString(`REPLACE INTO ` + s.TableName(db, table) + ` `)
	if len(columns) > 0 {
		stmt.WriteByte('(')
		for i, col := range columns {
			if i > 0 {
				stmt.WriteByte(',')
			}
			stmt.WriteString(s.Quote(col))
		}
		stmt.WriteString(`) `)
	}
	err = s.parser.BuildStatement(stmt, query)
	if err != nil {
		return
	}
	stmt.WriteByte(';')
	return
}
