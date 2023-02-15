package mysql

import (
	"regexp"
	"strconv"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql"
)

// HasIndexByName :
func (s *mySQL) HasIndexByName(stmt db.Stmt, dbName, table, indexName string) {
	stmt.WriteString(`SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?;`)
	stmt.AppendArgs(dbName, table, indexName)
}

// HasIndex :
func (s *mySQL) HasIndex(stmt db.Stmt, dbName, table string, idx sql.Index) {
	nonUnique, idxType := true, "BTREE"
	switch idx.Type {
	case sql.Unique:
		nonUnique = false
	case sql.FullText:
		idxType = "FULLTEXT"
	case sql.Spatial:
		idxType = "SPATIAL"
	case sql.Primary:
		nonUnique = false
	}
	args := []any{dbName, table, idxType, nonUnique}
	stmt.WriteString("SELECT COUNT(1) FROM (")
	stmt.WriteString("SELECT INDEX_NAME, COUNT(*) AS c FROM INFORMATION_SCHEMA.STATISTICS ")
	stmt.WriteString("WHERE TABLE_SCHEMA = ? ")
	stmt.WriteString("AND TABLE_NAME = ? ")
	stmt.WriteString("AND INDEX_TYPE = ? ")
	stmt.WriteString("AND NON_UNIQUE = ? ")
	stmt.WriteString(`AND COLUMN_NAME IN (`)
	for i, col := range idx.Columns {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteByte('?')
		args = append(args, col.Name)
	}
	stmt.WriteString(`) GROUP BY INDEX_NAME`)
	stmt.WriteString(") AS temp WHERE temp.c = ?")
	stmt.WriteByte(';')
	args = append(args, int64(len(idx.Columns)))
	stmt.AppendArgs(args...)
}

// GetIndexes :
func (s *mySQL) GetIndexes(stmt db.Stmt, dbName, table string) {
	stmt.WriteString(`SELECT DISTINCT INDEX_NAME, INDEX_TYPE, NON_UNIQUE FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;`)
	stmt.AppendArgs(dbName, table)
}

// CreateIndexes :
func (s *mySQL) CreateIndexes(stmt db.Stmt, db, table string, idxs []sql.Index, supportDesc bool) {
	stmt.WriteString(`ALTER TABLE ` + s.TableName(db, table))
	for i, idx := range idxs {
		if i > 0 {
			stmt.WriteByte(',')
		}

		stmt.WriteString(` ADD ` + s.getIndexByType(idx.Type) + ` `)
		name := idx.GetName()
		if idx.Type == sql.MultiValued {
			stmt.WriteString(name + "( (CAST(")
			if regexp.MustCompile(`(?is).+\s*\-\>\s*.+`).MatchString(idx.Cast) {
				stmt.WriteString(idx.Cast)
			} else {
				stmt.WriteString("`" + idx.Cast + "` -> '$'")
			}
			stmt.WriteString(` AS ` + idx.As + `)) )`)
		} else {
			if name != "" {
				stmt.WriteString(s.Quote(name))
			}
			stmt.WriteString(` (`)
			for j, col := range idx.Columns {
				if j > 0 {
					stmt.WriteByte(',')
				}
				stmt.WriteString(s.Quote(col.Name))
				if !supportDesc {
					continue
				}
				if col.Direction == sql.Descending {
					stmt.WriteString(` DESC`)
				}
			}
			stmt.WriteByte(')')
		}

		if idx.Comment != "" {
			stmt.WriteString(` COMMENT ` + strconv.Quote(idx.Comment))
		}
	}
	stmt.WriteByte(';')
}

// DropIndexes :
func (s *mySQL) DropIndexes(stmt db.Stmt, db, table string, idxs []string) {
	stmt.WriteString(`ALTER TABLE ` + s.TableName(db, table) + ` `)
	for i, idx := range idxs {
		if idx == "PRIMARY" {
			// stmt.WriteString("DROP PRIMARY KEY")
			continue
		}
		if i > 0 {
			stmt.WriteByte(',')
		}

		stmt.WriteString(`DROP INDEX ` + s.Quote(idx))
	}
	stmt.WriteByte(';')
}

func (s *mySQL) getIndexByType(k sql.Type) (idx string) {
	switch k {
	case sql.FullText:
		idx = "FULLTEXT INDEX"
	case sql.Spatial:
		idx = "SPATIAL INDEX"
	case sql.Unique:
		idx = "UNIQUE INDEX"
	case sql.Primary:
		idx = "PRIMARY KEY"
	default:
		idx = "INDEX"
	}
	return
}
