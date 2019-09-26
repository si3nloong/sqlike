package mysql

import (
	"reflect"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/charset"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sql/util"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/columns"
)

// RenameTable :
func (ms MySQL) RenameTable(db, oldName, newName string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`RENAME TABLE `)
	stmt.WriteString(ms.TableName(db, oldName))
	stmt.WriteString(` TO `)
	stmt.WriteString(ms.TableName(db, newName))
	stmt.WriteByte(';')
	return
}

// DropTable :
func (ms MySQL) DropTable(db, table string, exists bool) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("DROP TABLE")
	if exists {
		stmt.WriteString(" IF EXISTS")
	}
	stmt.WriteByte(' ')
	stmt.WriteString(ms.TableName(db, table) + ";")
	return
}

// TruncateTable :
func (ms MySQL) TruncateTable(db, table string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("TRUNCATE TABLE " + ms.TableName(db, table) + ";")
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
func (ms MySQL) CreateTable(db, table string, code charset.Code, collate, pk string, fields []*reflext.StructField) (stmt *sqlstmt.Statement, err error) {
	var (
		col     columns.Column
		virtual bool
		stored  bool
	)

	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`CREATE TABLE ` + ms.TableName(db, table) + ` `)
	stmt.WriteRune('(')

	// Main columns :
	for i, sf := range fields {
		if i > 0 {
			stmt.WriteRune(',')
		}

		col, err = ms.schema.GetColumn(sf)
		if err != nil {
			return
		}
		if sf.Path == pk {
			stmt.WriteString("PRIMARY KEY (`" + pk + "`)")
			stmt.WriteRune(',')
		}
		if _, isOk := sf.Tag.LookUp("unique_index"); isOk {
			stmt.WriteString("UNIQUE INDEX `UX_" + sf.Path + "`(`" + sf.Path + "`)")
			stmt.WriteRune(',')
		}

		ms.buildSchemaByColumn(stmt, col)

		// Generated columns :
		t := reflext.Deref(sf.Zero.Type())
		if t.Kind() != reflect.Struct {
			continue
		}

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
	stmt.WriteRune(')')
	stmt.WriteString(" ENGINE=INNODB")
	if code == "" {
		stmt.WriteString(" CHARACTER SET utf8mb4")
		stmt.WriteString(" COLLATE utf8mb4_unicode_ci")
	} else {
		stmt.WriteString(" CHARACTER SET " + string(code))
		if collate != "" {
			stmt.WriteString(" COLLATE " + collate)
		}
	}
	stmt.WriteRune(';')
	return
}

// AlterTable :
func (ms *MySQL) AlterTable(db, table, pk string, fields []*reflext.StructField, cols util.StringSlice, indexes util.StringSlice, unsafe bool) (stmt *sqlstmt.Statement, err error) {
	var (
		col     columns.Column
		idx     int
		virtual bool
		stored  bool
	)

	suffix := "FIRST"
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString(`ALTER TABLE ` + ms.TableName(db, table) + ` `)

	// TODO: add primary key when missing

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
		if action == "ADD" && sf.Path == pk {
			stmt.WriteString("ADD PRIMARY KEY (`" + pk + "`)")
			stmt.WriteRune(',')
		}
		stmt.WriteString(action + " ")
		col, err = ms.schema.GetColumn(sf)
		if err != nil {
			return
		}
		ms.buildSchemaByColumn(stmt, col)
		stmt.WriteString(" " + suffix)
		suffix = "AFTER " + ms.Quote(sf.Path)

		// Generated columns :
		t := reflext.Deref(sf.Zero.Type())
		if t.Kind() != reflect.Struct {
			continue
		}

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

				action = "ADD"
				idx = cols.IndexOf(child.Path)
				if idx > -1 {
					action = "MODIFY"
					cols.Splice(idx)
				}

				stmt.WriteString(action + " ")
				stmt.WriteString(ms.Quote(col.Name))
				stmt.WriteString(" " + col.Type)
				path := strings.TrimLeft(strings.TrimPrefix(child.Path, sf.Path), ".")
				stmt.WriteString(" AS ")
				stmt.WriteString("(" + ms.Quote(sf.Path) + "->>'$." + path + "')")
				if stored {
					stmt.WriteString(" STORED")
				}
				if !col.Nullable {
					stmt.WriteString(" NOT NULL")
				}
				suffix = "AFTER " + ms.Quote(child.Path)
			}
			children = children[1:]
			children = append(children, child.Children...)
		}
	}

	if unsafe {
		for _, col := range cols {
			stmt.WriteByte(',')
			stmt.WriteString("DROP COLUMN ")
			stmt.WriteString(ms.Quote(col))
		}
	}
	// stmt.WriteRune(',')
	// stmt.WriteString(`CONVERT TO CHARACTER SET utf8mb4`)
	// stmt.WriteString(` COLLATE utf8mb4_unicode_ci`)
	stmt.WriteRune(';')
	return
}

// Copy :
func (ms MySQL) Copy(db, table string, columns []string, act *actions.CopyActions) (stmt *sqlstmt.Statement, err error) {
	stmt = new(sqlstmt.Statement)
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
	err = ms.parser.BuildStatement(stmt, &act.FindActions)
	return
}

func (ms MySQL) buildSchemaByColumn(stmt *sqlstmt.Statement, col columns.Column) {
	stmt.WriteString(ms.Quote(col.Name))
	stmt.WriteString(" " + col.Type)
	if col.CharSet != nil {
		stmt.WriteString(" CHARACTER SET " + *col.CharSet)
	}
	if col.Collation != nil {
		stmt.WriteString(" COLLATE " + *col.Collation)
	}
	if col.Extra != "" {
		stmt.WriteString(" " + col.Extra)
	}
	if !col.Nullable {
		stmt.WriteString(" NOT NULL")
		if col.DefaultValue != nil {
			stmt.WriteString(" DEFAULT " + ms.WrapOnlyValue(*col.DefaultValue))
		}
	}
}
