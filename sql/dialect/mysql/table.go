package mysql

import (
	"reflect"
	"strings"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/sql/driver"
	"github.com/si3nloong/sqlike/v2/sql/util"
	"github.com/si3nloong/sqlike/v2/sqlike/indexes"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// HasPrimaryKey :
func (ms MySQL) HasPrimaryKey(stmt db.Stmt, db, table string) {
	stmt.WriteString("SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS ")
	stmt.WriteString("WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND CONSTRAINT_TYPE = 'PRIMARY KEY'")
	stmt.WriteByte(';')
	stmt.AppendArgs(db, table)
}

// RenameTable :
func (ms MySQL) RenameTable(stmt db.Stmt, db, oldName, newName string) {
	stmt.WriteString("RENAME TABLE ")
	stmt.WriteString(ms.TableName(db, oldName))
	stmt.WriteString(" TO ")
	stmt.WriteString(ms.TableName(db, newName))
	stmt.WriteByte(';')
}

// DropTable :
func (ms MySQL) DropTable(stmt db.Stmt, db, table string, exists bool) {
	stmt.WriteString("DROP TABLE")
	if exists {
		stmt.WriteString(" IF EXISTS")
	}
	stmt.WriteByte(' ')
	stmt.WriteString(ms.TableName(db, table) + ";")
}

// TruncateTable :
func (ms MySQL) TruncateTable(stmt db.Stmt, db, table string) {
	stmt.WriteString("TRUNCATE TABLE " + ms.TableName(db, table) + ";")
}

// HasTable :
func (ms MySQL) HasTable(stmt db.Stmt, dbName, table string) {
	stmt.WriteString(`SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;`)
	stmt.AppendArgs(dbName, table)
}

// CreateTable :
func (ms MySQL) CreateTable(
	stmt db.Stmt,
	dbName, table, pk string,
	info driver.Info,
	fields []reflext.StructFielder,
) (err error) {
	var (
		col     *sql.Column
		pkk     reflext.StructFielder
		k1, k2  string
		virtual bool
		stored  bool
		ctx     = sql.Context(dbName, table)
	)

	stmt.WriteString("CREATE TABLE " + ms.TableName(dbName, table) + " ")
	stmt.WriteByte('(')

	// Main columns :
	for i, f := range fields {
		if i > 0 {
			stmt.WriteByte(',')
		}

		ctx.SetField(f)
		col, err = ms.schema.GetColumn(ctx)
		if err != nil {
			return
		}

		tag := f.Tag()
		// allow primary_key tag to override
		if _, ok := tag.LookUp("primary_key"); ok {
			pkk = f
		} else if _, ok := tag.LookUp("auto_increment"); ok {
			pkk = f
		} else if f.Name() == pk && pkk == nil {
			pkk = f
		}

		idx := indexes.Index{Columns: indexes.Columns(f.Name())}
		if _, ok := tag.LookUp("unique_index"); ok {
			stmt.WriteString("UNIQUE INDEX " + idx.GetName() + " (" + ms.Quote(f.Name()) + ")")
			stmt.WriteByte(',')
		}

		ms.buildSchemaByColumn(stmt, col)

		if v, ok := tag.LookUp("comment"); ok {
			if len(v) > 60 {
				panic("sqlike: maximum length of comment is 60 characters")
			}
			stmt.WriteString(" COMMENT '" + v + "'")
		}

		// check generated columns
		t := reflext.Deref(f.Type())
		if t.Kind() != reflect.Struct {
			continue
		}

		children := f.Children()
		for len(children) > 0 {
			child := children[0]
			tg := child.Tag()
			k1, virtual = tg.LookUp("virtual_column")
			k2, stored = tg.LookUp("stored_column")
			if virtual || stored {
				stmt.WriteByte(',')

				ctx.SetField(child)
				col, err = ms.schema.GetColumn(ctx)
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
				path := strings.TrimLeft(strings.TrimPrefix(child.Name(), f.Name()), ".")
				stmt.WriteString(" AS ")
				stmt.WriteString("(" + ms.Quote(f.Name()) + "->>'$." + path + "')")
				if stored {
					stmt.WriteString(" STORED")
				}
				if !col.Nullable {
					stmt.WriteString(" NOT NULL")
				}
			}
			children = children[1:]
			children = append(children, child.Children()...)
		}

	}
	if pkk != nil {
		stmt.WriteByte(',')
		stmt.WriteString("PRIMARY KEY (" + ms.Quote(pkk.Name()) + ")")
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
func (ms *MySQL) AlterTable(
	stmt db.Stmt,
	dbName, table, pk string,
	hasPk bool,
	info driver.Info,
	fields []reflext.StructFielder,
	cols util.StringSlice,
	idxs util.StringSlice,
	unsafe bool,
) (err error) {
	var (
		col     *sql.Column
		pkk     reflext.StructFielder
		idx     int
		k1, k2  string
		virtual bool
		stored  bool
		ctx     = sql.Context(dbName, table)
	)

	suffix := "FIRST"
	stmt.WriteString("ALTER TABLE " + ms.TableName(dbName, table) + " ")

	for i, f := range fields {
		if i > 0 {
			stmt.WriteByte(',')
		}

		action := "ADD"
		idx = cols.IndexOf(f.Name())
		if idx > -1 {
			action = "MODIFY"
			cols.Splice(idx)
		}
		if !hasPk {
			// allow primary_key tag to override
			if _, ok := f.Tag().LookUp("primary_key"); ok {
				pkk = f
			}
			if f.Name() == pk && pkk == nil {
				pkk = f
			}
		}

		tag := f.Tag()
		_, ok1 := tag.LookUp("unique_index")
		_, ok2 := tag.LookUp("auto_increment")
		if ok1 || ok2 {
			idx := indexes.Index{Columns: indexes.Columns(f.Name())}
			if idxs.IndexOf(idx.GetName()) < 0 {
				stmt.WriteString("ADD")
				stmt.WriteString(" UNIQUE INDEX " + idx.GetName() + " (" + ms.Quote(f.Name()) + ")")
				stmt.WriteByte(',')
			}
		}
		stmt.WriteString(action + " ")

		ctx.SetField(f)
		col, err = ms.schema.GetColumn(ctx)
		if err != nil {
			return
		}

		ms.buildSchemaByColumn(stmt, col)

		if v, ok := f.Tag().LookUp("comment"); ok {
			if len(v) > 60 {
				panic("maximum length of comment is 60 characters")
			}
			stmt.WriteString(" COMMENT '" + v + "'")
		}

		stmt.WriteString(" " + suffix)
		suffix = "AFTER " + ms.Quote(f.Name())

		// check generated columns
		t := reflext.Deref(f.Type())
		if t.Kind() != reflect.Struct {
			continue
		}

		children := f.Children()
		for len(children) > 0 {
			child := children[0]
			tg := child.Tag()
			k1, virtual = tg.LookUp("virtual_column")
			k2, stored = tg.LookUp("stored_column")
			if virtual || stored {
				stmt.WriteByte(',')
				ctx.SetField(child)
				col, err = ms.schema.GetColumn(ctx)
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
				path := strings.TrimLeft(strings.TrimPrefix(child.Name(), f.Name()), ".")
				stmt.WriteString(" AS ")
				stmt.WriteString("(" + ms.Quote(f.Name()) + "->>'$." + path + "')")
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
			children = append(children, child.Children()...)
		}

	}

	if pkk != nil {
		stmt.WriteByte(',')
		stmt.WriteString("ADD PRIMARY KEY (" + ms.Quote(pkk.Name()) + ")")
	}

	if unsafe {
		for _, col := range cols {
			stmt.WriteByte(',')
			stmt.WriteString("DROP COLUMN ")
			stmt.WriteString(ms.Quote(col))
		}
	}

	// TODO: character set
	stmt.WriteByte(',')
	stmt.WriteString(`CHARACTER SET utf8mb4`)
	stmt.WriteString(` COLLATE utf8mb4_unicode_ci`)
	stmt.WriteByte(';')
	return
}
