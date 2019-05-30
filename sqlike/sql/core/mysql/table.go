package mysql

import (
	"log"
	"reflect"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/sql/component"
	sqlstmt "github.com/si3nloong/sqlike/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/sql/util"
)

// RenameTable :
func (ms MySQL) RenameTable(oldName, newName string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`RENAME TABLE `)
	stmt.WriteString(ms.Quote(oldName))
	stmt.WriteString(` TO `)
	stmt.WriteString(ms.Quote(newName))
	stmt.WriteByte(';')
	return
}

// DropTable :
func (ms MySQL) DropTable(table string, exists bool) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`DROP TABLE`)
	if exists {
		stmt.WriteString(` IF EXISTS`)
	}
	stmt.WriteString(` ` + ms.Quote(table) + `;`)
	return
}

// TruncateTable :
func (ms MySQL) TruncateTable(table string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`TRUNCATE TABLE ` + ms.Quote(table) + `;`)
	return
}

// HasTable :
func (ms MySQL) HasTable(dbName, table string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;`)
	stmt.AppendArgs([]interface{}{dbName, table})
	return
}

// CreateTable :
func (ms MySQL) CreateTable(table string, fields []*reflext.StructField) (stmt *sqlstmt.Statement, err error) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`CREATE TABLE ` + ms.Quote(table) + ` `)
	stmt.WriteRune('(')

	var col component.Column
	var virtual, stored bool
	// Main columns :
	for i, sf := range fields {
		if i > 0 {
			stmt.WriteRune(',')
		}

		col, err = ms.schema.GetColumn(sf)
		if err != nil {
			return
		}
		if sf.Path == "$Key" {
			stmt.WriteString("PRIMARY KEY (`$Key`)")
			stmt.WriteRune(',')
		}
		ms.buildSchemaByColumn(stmt, col)

		// Generated columns :
		t := reflext.Deref(sf.Zero.Type())
		if t.Kind() == reflect.Struct {
			children := sf.Children
			for len(children) > 0 {
				child := children[0]
				_, virtual = child.Tag.LookUp("virtual_column")
				_, stored = child.Tag.LookUp("stored_column")
				if virtual || stored {
					stmt.WriteRune(',')
					col, err = ms.schema.GetColumn(child)
					if err != nil {
						return
					}

					stmt.WriteString(ms.Quote(col.Name))
					stmt.WriteString(` ` + col.Type)
					path := strings.TrimLeft(strings.TrimPrefix(child.Path, sf.Path), `.`)
					stmt.WriteString(` AS `)
					stmt.WriteString(`(` + ms.Quote(sf.Path) + `->>'$.` + path + `')`)
					if stored {
						stmt.WriteString(` STORED`)
					}
					if !col.Nullable {
						stmt.WriteString(` NOT NULL`)
					}
				}
				children = children[1:]
				children = append(children, child.Children...)
			}
		}
	}
	stmt.WriteRune(')')
	stmt.WriteString(` ENGINE=INNODB`)
	stmt.WriteString(` CHARACTER SET utf8mb4`)
	stmt.WriteString(` COLLATE utf8mb4_unicode_ci`)
	stmt.WriteRune(';')
	return
}

// AlterTable :
func (ms *MySQL) AlterTable(table string, fields []*reflext.StructField, columns []string, indexes []string, unsafe bool) (stmt *sqlstmt.Statement, err error) {
	stmt = sqlstmt.NewStatement(ms)
	cols := util.StringSlice(columns)
	stmt.WriteString(`ALTER TABLE ` + ms.Quote(table) + ` `)

	var col component.Column
	var idx int
	for i, sf := range fields {
		if i > 0 {
			stmt.WriteRune(',')
		}
		action := "ADD"
		idx = cols.IndexOf(sf.Path)
		if idx > -1 {
			action = "MODIFY"
			cols.Splice(idx)
		}

		stmt.WriteString(action + ` `)
		col, err = ms.schema.GetColumn(sf)
		if err != nil {
			return
		}
		ms.buildSchemaByColumn(stmt, col)
		if i < 1 {
			stmt.WriteString(` FIRST`)
		} else {
			stmt.WriteString(` AFTER ` + ms.Quote(fields[i-1].Path))
		}

		if action == "ADD" && sf.Path == "$Key" {
			stmt.WriteRune(',')
			stmt.WriteString("ADD PRIMARY KEY (`$Key`)")
		}
	}

	if unsafe {
		// stmt.WriteString(` DROP`)
		for _, col := range cols {
			log.Println(col)
		}
		log.Println("Unsafe columns :")
	}
	// stmt.WriteRune(',')
	// stmt.WriteString(`CONVERT TO CHARACTER SET utf8mb4`)
	// stmt.WriteString(` COLLATE utf8mb4_unicode_ci`)
	stmt.WriteRune(';')
	return
}

func (ms MySQL) buildSchemaByColumn(stmt *sqlstmt.Statement, col component.Column) {
	stmt.WriteString(ms.Quote(col.Name))
	stmt.WriteString(` ` + col.Type)
	if col.CharSet != nil {
		stmt.WriteString(` CHARACTER SET ` + *col.CharSet)
	}
	if col.Collation != nil {
		stmt.WriteString(` COLLATE ` + *col.Collation)
	}
	if col.Extra != "" {
		stmt.WriteString(` ` + col.Extra)
	}
	if !col.Nullable {
		stmt.WriteString(` NOT NULL`)
		if col.DefaultValue != nil {
			stmt.WriteString(` DEFAULT ` + ms.WrapOnlyValue(*col.DefaultValue))
		}
	}
}
