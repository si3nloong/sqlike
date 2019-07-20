package mysql

import (
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/indexes"
)

func (ms MySQL) HasIndex(dbName, table, indexName string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?;`)
	stmt.AppendArgs([]interface{}{dbName, table, indexName})
	return
}

// GetIndexes :
func (ms MySQL) GetIndexes(dbName, table string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`SELECT DISTINCT INDEX_NAME, INDEX_TYPE, IS_VISIBLE FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;`)
	stmt.AppendArgs([]interface{}{dbName, table})
	return
}

// CreateIndexes :
func (ms MySQL) CreateIndexes(table string, idxs []indexes.Index, supportDesc bool) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`ALTER TABLE ` + ms.Quote(table))
	for i, idx := range idxs {
		if i > 0 {
			stmt.WriteRune(',')
		}

		stmt.WriteString(` ADD ` + ms.getIndexByType(idx.Type))
		stmt.WriteString(` ` + ms.Quote(idx.GetName()) + ` `)
		stmt.WriteRune('(')
		for j, col := range idx.Columns {
			if j > 0 {
				stmt.WriteRune(',')
			}
			stmt.WriteString(ms.Quote(col.Name))
			if !supportDesc {
				continue
			}
			if col.Direction == indexes.Descending {
				stmt.WriteString(" DESC")
			}
		}
		stmt.WriteRune(')')
	}
	stmt.WriteRune(';')
	return
}

// DropIndex :
func (ms MySQL) DropIndex(table, idxName string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`DROP INDEX`)
	stmt.WriteString(` ` + ms.Quote(idxName))
	stmt.WriteString(` ON ` + ms.Quote(table))
	stmt.WriteRune(';')
	return
}

func (ms MySQL) getIndexByType(k indexes.Type) (idx string) {
	switch k {
	case indexes.FullText:
		idx = `FULLTEXT `
	case indexes.Spatial:
		idx = `SPATIAL `
	case indexes.Unique:
		idx = `UNIQUE `
	}
	idx += `INDEX`
	return
}
