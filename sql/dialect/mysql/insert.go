package mysql

import (
	"reflect"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// InsertInto :
func (ms *mySQL) InsertInto(
	stmt db.Stmt,
	dbName, table, pk string,
	cache reflext.StructMapper,
	fields []reflext.FieldInfo,
	v reflect.Value,
	opt *options.InsertOptions,
) (err error) {
	noOfRecords := v.Len()

	stmt.WriteString("INSERT")
	if opt.Mode == options.InsertIgnore {
		stmt.WriteString(" IGNORE")
	}
	stmt.WriteString(" INTO " + ms.TableName(dbName, table) + " (")

	omitField := make(map[string]bool)
	noOfOmit := len(opt.Omits)
	for i := 0; i < len(fields); {
		// omit all the field provided by user
		if noOfOmit > 0 && opt.Omits.IndexOf(fields[i].Name()) > -1 {
			if opt.Mode != options.InsertOnDuplicate {
				fields = append(fields[:i], fields[i+1:]...)
				continue
			} else {
				omitField[fields[i].Name()] = true
			}
		}

		// omit all the struct field with `generated_column` tag, it shouldn't include when inserting to the db
		if _, ok := fields[i].Tag().Option("generated_column"); ok {
			fields = append(fields[:i], fields[i+1:]...)
			continue
		}

		stmt.WriteString(ms.Quote(fields[i].Name()))
		if i < len(fields)-1 {
			stmt.WriteByte(',')
		}

		i++
	}
	stmt.WriteString(`) VALUES `)

	length := len(fields)
	encoders := make([]db.ValueEncoder, length)
	for i := 0; i < noOfRecords; i++ {
		if i > 0 {
			stmt.WriteByte(',')
		}

		vi := reflext.Indirect(v.Index(i))
		stmt.WriteByte('(')
		// marshal records and construct `VALUES` statement
		for j, f := range fields {
			if j > 0 {
				stmt.WriteByte(',')
			}

			// get struct property value
			fv := cache.FieldByIndexesReadOnly(vi, f.Index())

			// first record only find encoders
			if i == 0 {
				encoders[j], err = ms.LookupEncoder(fv)
				if err != nil {
					return err
				}
			}

			query, args, err := encoders[j](ms, fv, f.Tag().Opts())
			if err != nil {
				return err
			}
			stmt.AppendArgs(query, args...)
		}
		stmt.WriteByte(')')
	}

	if opt.Mode == options.InsertOnDuplicate {
		var (
			column string
			name   string
		)
		stmt.WriteString(` ON DUPLICATE KEY UPDATE `)
		next := false
		for _, f := range fields {
			name = f.Name()
			// skip primary key on duplicate update
			if name == pk {
				continue
			}

			// skip primary key on duplicate update
			if _, ok := f.Tag().Option("primary_key"); ok {
				continue
			}

			if _, ok := f.Tag().Option("auto_increment"); ok {
				continue
			}

			// skip omit fields on update
			if _, ok := omitField[name]; ok {
				continue
			}

			if next {
				stmt.WriteByte(',')
			}

			column = ms.Quote(name)
			stmt.WriteString(column + `=VALUES(` + column + `)`)
			next = true
		}
	}
	stmt.WriteByte(';')
	return
}
