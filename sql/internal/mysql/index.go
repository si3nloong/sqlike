package mysql

import (
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/indexes"
)

// HasIndexByName :
func (ms MySQL) HasIndexByName(stmt sqlstmt.Stmt, dbName, table, indexName string) {
	stmt.WriteString(`SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?;`)
	stmt.AppendArgs(dbName, table, indexName)
}

// HasIndex :
func (ms MySQL) HasIndex(stmt sqlstmt.Stmt, dbName, table string, idx indexes.Index) {
	nonUnique, idxType := true, "BTREE"
	switch idx.Type {
	case indexes.Unique:
		nonUnique = false
	case indexes.FullText:
		idxType = "FULLTEXT"
	case indexes.Spatial:
		idxType = "SPATIAL"
	case indexes.Primary:
		nonUnique = false
	}
	args := []interface{}{dbName, table, idxType, nonUnique}
	stmt.WriteString("SELECT COUNT(1) FROM (")
	stmt.WriteString("SELECT INDEX_NAME, COUNT(*) AS c FROM INFORMATION_SCHEMA.STATISTICS ")
	stmt.WriteString("WHERE TABLE_SCHEMA = ? ")
	stmt.WriteString("AND TABLE_NAME = ? ")
	stmt.WriteString("AND INDEX_TYPE = ? ")
	stmt.WriteString("AND NON_UNIQUE = ? ")
	stmt.WriteString("AND COLUMN_NAME IN ")
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
	args = append(args, int64(len(idx.Columns)))
	stmt.AppendArgs(args...)
}

// GetIndexes :
func (ms MySQL) GetIndexes(stmt sqlstmt.Stmt, dbName, table string) {
	stmt.WriteString(`SELECT DISTINCT INDEX_NAME, INDEX_TYPE, NON_UNIQUE FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;`)
	stmt.AppendArgs(dbName, table)
}

// CreateIndexes :
func (ms MySQL) CreateIndexes(stmt sqlstmt.Stmt, db, table string, idxs []indexes.Index, supportDesc bool) {
	stmt.WriteString("ALTER TABLE " + ms.TableName(db, table))
	for i, idx := range idxs {
		if i > 0 {
			stmt.WriteByte(',')
		}

		stmt.WriteString(" ADD " + ms.getIndexByType(idx.Type))
		name := idx.GetName()
		if name != "" {
			stmt.WriteString(" " + ms.Quote(name))
		}
		stmt.WriteString(" (")
		for j, col := range idx.Columns {
			if j > 0 {
				stmt.WriteByte(',')
			}
			stmt.WriteString(ms.Quote(col.Name))
			if !supportDesc {
				continue
			}
			if col.Direction == indexes.Descending {
				stmt.WriteString(" DESC")
			}
		}
		stmt.WriteByte(')')
	}
	stmt.WriteByte(';')
}

// DropIndex :
func (ms MySQL) DropIndexes(stmt sqlstmt.Stmt, db, table string, idxs []string) {
	stmt.WriteString("ALTER TABLE " + ms.TableName(db, table) + " ")
	for i, idx := range idxs {
		if idx == "PRIMARY" {
			// stmt.WriteString("DROP PRIMARY KEY")
			continue
		}
		if i > 0 {
			stmt.WriteByte(',')
		}

		stmt.WriteString("DROP INDEX ")
		stmt.WriteString(ms.Quote(idx))
	}
	stmt.WriteByte(';')
}

func (ms MySQL) getIndexByType(k indexes.Type) (idx string) {
	switch k {
	case indexes.FullText:
		idx = "FULLTEXT INDEX"
	case indexes.Spatial:
		idx = "SPATIAL INDEX"
	case indexes.Unique:
		idx = "UNIQUE INDEX"
	case indexes.Primary:
		idx = "PRIMARY KEY"
	default:
		idx = "INDEX"
	}
	return
}
