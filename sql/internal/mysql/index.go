package mysql

import (
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/indexes"
)

// HasIndexByName :
func (ms MySQL) HasIndexByName(dbName, table, indexName string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?;`)
	stmt.AppendArgs([]interface{}{dbName, table, indexName})
	return
}

// HasIndex :
func (ms MySQL) HasIndex(dbName, table string, idx indexes.Index) (stmt *sqlstmt.Statement) {
	nonUnique, idxType := true, "BTREE"
	switch idx.Type {
	case indexes.Unique:
		nonUnique = false
	case indexes.FullText:
		idxType = "FULLTEXT"
	case indexes.Spatial:
		idxType = "SPATIAL"
	}
	args := []interface{}{dbName, table, idxType, nonUnique}
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("SELECT COUNT(1) FROM (")
	stmt.WriteString("SELECT INDEX_NAME, COUNT(*) AS c")
	stmt.WriteString(" FROM INFORMATION_SCHEMA.STATISTICS")
	stmt.WriteString(" WHERE TABLE_SCHEMA = ?")
	stmt.WriteString(" AND TABLE_NAME = ?")
	stmt.WriteString(" AND INDEX_TYPE = ?")
	stmt.WriteString(" AND NON_UNIQUE = ?")
	stmt.WriteString(" AND COLUMN_NAME IN ")
	stmt.WriteByte('(')
	for i, col := range idx.Columns {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteByte('?')
		args = append(args, col.Name)
	}
	stmt.WriteByte(')')
	stmt.WriteString(" GROUP BY INDEX_NAME")
	stmt.WriteString(") AS temp WHERE temp.c = ?")
	stmt.WriteByte(';')
	stmt.AppendArgs(append(args, int64(len(idx.Columns))))
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
func (ms MySQL) CreateIndexes(db, table string, idxs []indexes.Index, supportDesc bool) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`ALTER TABLE ` + ms.TableName(db, table))
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
func (ms MySQL) DropIndex(db, table, idxName string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`DROP INDEX`)
	stmt.WriteString(` ` + ms.Quote(idxName))
	stmt.WriteString(` ON ` + ms.TableName(db, table))
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
