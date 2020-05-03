package mysql

import (
	"reflect"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sql/util"
	"github.com/si3nloong/sqlike/sqlike/columns"
	"github.com/si3nloong/sqlike/sqlike/indexes"
)

// HasPrimaryKey :
func (ms MySQL) HasPrimaryKey(stmt sqlstmt.Stmt, db, table string) {
	stmt.WriteString("SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS ")
	stmt.WriteString("WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND CONSTRAINT_TYPE = 'PRIMARY KEY'")
	stmt.WriteByte(';')
	stmt.AppendArgs(db, table)
}

// RenameTable :
func (ms MySQL) RenameTable(stmt sqlstmt.Stmt, db, oldName, newName string) {
	stmt.WriteString("RENAME TABLE ")
	stmt.WriteString(ms.TableName(db, oldName))
	stmt.WriteString(" TO ")
	stmt.WriteString(ms.TableName(db, newName))
	stmt.WriteByte(';')
}

// DropTable :
func (ms MySQL) DropTable(stmt sqlstmt.Stmt, db, table string, exists bool) {
	stmt.WriteString("DROP TABLE")
	if exists {
		stmt.WriteString(" IF EXISTS")
	}
	stmt.WriteByte(' ')
	stmt.WriteString(ms.TableName(db, table) + ";")
}

// TruncateTable :
func (ms MySQL) TruncateTable(stmt sqlstmt.Stmt, db, table string) {
	stmt.WriteString("TRUNCATE TABLE " + ms.TableName(db, table) + ";")
}

// HasTable :
func (ms MySQL) HasTable(stmt sqlstmt.Stmt, dbName, table string) {
	stmt.WriteString(`SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;`)
	stmt.AppendArgs(dbName, table)
}

// CreateTable :
func (ms MySQL) CreateTable(stmt sqlstmt.Stmt, db, table, pk string, info driver.Info, fields []*reflext.StructField) (err error) {
	var (
		col     columns.Column
		pkk     *reflext.StructField
		k1, k2  string
		virtual bool
		stored  bool
	)

	stmt.WriteString("CREATE TABLE " + ms.TableName(db, table) + " ")
	stmt.WriteByte('(')

	// Main columns :
	for i, sf := range fields {
		if i > 0 {
			stmt.WriteByte(',')
		}

		col, err = ms.schema.GetColumn(info, sf)
		if err != nil {
			return
		}

		// allow primary_key tag to override
		if _, ok := sf.Tag.LookUp("primary_key"); ok {
			pkk = sf
		} else if _, ok := sf.Tag.LookUp("auto_increment"); ok {
			pkk = sf
		} else if sf.Path == pk && pkk == nil {
			pkk = sf
		}

		idx := indexes.Index{Columns: indexes.Columns(sf.Path)}
		if _, ok := sf.Tag.LookUp("unique_index"); ok {
			stmt.WriteString("UNIQUE INDEX " + idx.GetName() + " (" + ms.Quote(sf.Path) + ")")
			stmt.WriteByte(',')
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
				stmt.WriteByte(',')
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
	if pkk != nil {
		stmt.WriteByte(',')
		stmt.WriteString("PRIMARY KEY (" + ms.Quote(pkk.Path) + ")")
	}
	stmt.WriteByte(')')
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
	stmt.WriteByte(';')
	return
}

// AlterTable :
func (ms *MySQL) AlterTable(stmt sqlstmt.Stmt, db, table, pk string, hasPk bool, info driver.Info, fields []*reflext.StructField, cols util.StringSlice, idxs util.StringSlice, unsafe bool) (err error) {
	var (
		col     columns.Column
		pkk     *reflext.StructField
		idx     int
		k1, k2  string
		virtual bool
		stored  bool
	)

	suffix := "FIRST"
	stmt.WriteString("ALTER TABLE " + ms.TableName(db, table) + " ")

	for i, sf := range fields {
		if i > 0 {
			stmt.WriteByte(',')
		}

		action := "ADD"
		idx = cols.IndexOf(sf.Path)
		if idx > -1 {
			action = "MODIFY"
			cols.Splice(idx)
		}
		if !hasPk {
			// allow primary_key tag to override
			if _, ok := sf.Tag.LookUp("primary_key"); ok {
				pkk = sf
			}
			if sf.Path == pk && pkk == nil {
				pkk = sf
			}
		}

		_, ok1 := sf.Tag.LookUp("unique_index")
		_, ok2 := sf.Tag.LookUp("auto_increment")
		if ok1 || ok2 {
			idx := indexes.Index{Columns: indexes.Columns(sf.Path)}
			if idxs.IndexOf(idx.GetName()) < 0 {
				stmt.WriteString("ADD")
				stmt.WriteString(" UNIQUE INDEX " + idx.GetName() + " (" + ms.Quote(sf.Path) + ")")
				stmt.WriteByte(',')
			}
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
				stmt.WriteByte(',')
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

	if pkk != nil {
		stmt.WriteByte(',')
		stmt.WriteString("ADD PRIMARY KEY (" + ms.Quote(pkk.Path) + ")")
	}

	if unsafe {
		for _, col := range cols {
			stmt.WriteByte(',')
			stmt.WriteString("DROP COLUMN ")
			stmt.WriteString(ms.Quote(col))
		}
	}

	// TODO: character set
	// stmt.WriteByte(',')
	// stmt.WriteString(`CONVERT TO CHARACTER SET utf8mb4`)
	// stmt.WriteString(` COLLATE utf8mb4_unicode_ci`)
	stmt.WriteByte(';')
	return
}
