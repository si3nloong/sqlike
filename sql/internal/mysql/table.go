package mysql

import (
	"reflect"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sql/util"
	"github.com/si3nloong/sqlike/sqlike/columns"
)

// RenameTable :
func (ms MySQL) RenameTable(db, oldName, newName string) (stmt *sqlstmt.Statement) {
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("RENAME TABLE ")
	stmt.WriteString(ms.TableName(db, oldName))
	stmt.WriteString(" TO ")
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
func (ms MySQL) CreateTable(db, table, pk string, info driver.Info, fields []*reflext.StructField) (stmt *sqlstmt.Statement, err error) {
	var (
		col     columns.Column
		k1, k2  string
		virtual bool
		stored  bool
	)

	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("CREATE TABLE " + ms.TableName(db, table) + " ")
	stmt.WriteRune('(')

	// Main columns :
	for i, sf := range fields {
		if i > 0 {
			stmt.WriteRune(',')
		}

		col, err = ms.schema.GetColumn(info, sf)
		if err != nil {
			return
		}
		_, ok := sf.Tag.LookUp("primary_key")
		if ok || sf.Path == pk {
			stmt.WriteString("PRIMARY KEY (" + ms.Quote(sf.Path) + ")")
			stmt.WriteRune(',')
		}
		if _, ok := sf.Tag.LookUp("unique_index"); ok {
			stmt.WriteString("UNIQUE INDEX " + ms.Quote("UX_"+sf.Path) + " (" + ms.Quote(sf.Path) + ")")
			stmt.WriteRune(',')
		}

		ms.buildSchemaByColumn(stmt, col)

		// Generated columns :
		t := reflext.Deref(sf.Type)
		if t.Kind() != reflect.Struct {
			continue
		}

		children := sf.Children
		for len(children) > 0 {
			child := children[0]
			k1, virtual = child.Tag.LookUp("virtual_column")
			k2, stored = child.Tag.LookUp("stored_column")
			if virtual || stored {
				stmt.WriteRune(',')
				col, err = ms.schema.GetColumn(info, child)
				if err != nil {
					return
				}

				name := col.Name
				if virtual && k1 != "" {
					name = k1
				}
				if stored && k2 != "" {
					name = k2
				}

				stmt.WriteString(ms.Quote(name))
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
			}
			children = children[1:]
			children = append(children, child.Children...)
		}
	}
	stmt.WriteRune(')')
	stmt.WriteString(" ENGINE=INNODB")
	code := string(info.Charset())
	if code == "" {
		stmt.WriteString(" CHARACTER SET utf8mb4")
		stmt.WriteString(" COLLATE utf8mb4_unicode_ci")
	} else {
		stmt.WriteString(" CHARACTER SET " + code)
		if info.Collate() != "" {
			stmt.WriteString(" COLLATE " + info.Collate())
		}
	}
	stmt.WriteRune(';')
	return
}

// AlterTable :
func (ms *MySQL) AlterTable(db, table, pk string, info driver.Info, fields []*reflext.StructField, cols util.StringSlice, indexes util.StringSlice, unsafe bool) (stmt *sqlstmt.Statement, err error) {
	var (
		col     columns.Column
		idx     int
		k1, k2  string
		virtual bool
		stored  bool
	)

	suffix := "FIRST"
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("ALTER TABLE " + ms.TableName(db, table) + " ")

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
		col, err = ms.schema.GetColumn(info, sf)
		if err != nil {
			return
		}
		ms.buildSchemaByColumn(stmt, col)
		stmt.WriteString(" " + suffix)
		suffix = "AFTER " + ms.Quote(sf.Path)

		// Generated columns :
		t := reflext.Deref(sf.Type)
		if t.Kind() != reflect.Struct {
			continue
		}

		children := sf.Children
		for len(children) > 0 {
			child := children[0]
			k1, virtual = child.Tag.LookUp("virtual_column")
			k2, stored = child.Tag.LookUp("stored_column")
			if virtual || stored {
				stmt.WriteRune(',')
				col, err = ms.schema.GetColumn(info, child)
				if err != nil {
					return
				}

				name := col.Name
				if virtual && k1 != "" {
					name = k1
				}
				if stored && k2 != "" {
					name = k2
				}

				action = "ADD"
				idx = cols.IndexOf(name)
				if idx > -1 {
					action = "MODIFY"
					cols.Splice(idx)
				}

				stmt.WriteString(action + " ")
				stmt.WriteString(ms.Quote(name))
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
				stmt.WriteString(" " + suffix)
				suffix = "AFTER " + ms.Quote(name)
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

	// TODO: character set
	// stmt.WriteRune(',')
	// stmt.WriteString(`CONVERT TO CHARACTER SET utf8mb4`)
	// stmt.WriteString(` COLLATE utf8mb4_unicode_ci`)
	stmt.WriteRune(';')
	return
}
