package mysql

import (
	"reflect"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/codec"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// InsertInto :
func (ms MySQL) InsertInto(db, table, pk string, mapper *reflext.Mapper, registry *codec.Registry, fields []*reflext.StructField, v reflect.Value, opt *options.InsertOptions) (stmt *sqlstmt.Statement, err error) {
	records := v.Len()
	stmt = sqlstmt.NewStatement(ms)
	stmt.WriteString("INSERT")
	if opt.Mode == options.InsertIgnore {
		stmt.WriteString(" IGNORE")
	}
	stmt.WriteString(" INTO " + ms.TableName(db, table) + " (")
	for i, f := range fields {
		if i > 0 {
			stmt.WriteRune(',')
		}
		stmt.WriteString(ms.Quote(f.Path))
	}
	stmt.WriteString(") VALUES ")
	length := len(fields)
	binds := strings.Repeat("?,", length)
	binds = "(" + binds[:len(binds)-1] + ")"
	encoders := make([]codec.ValueEncoder, length, length)
	for i := 0; i < records; i++ {
		if i > 0 {
			stmt.WriteRune(',')
		}
		vi := reflext.Indirect(v.Index(i))
		for j, sf := range fields {
			// first record only find encoders
			fv := mapper.FieldByIndexesReadOnly(vi, sf.Index)
			if i == 0 {
				encoders[j], err = findEncoder(mapper, registry, sf, fv)
				if err != nil {
					return nil, err
				}
			}

			val, err := encoders[j](sf, fv)
			if err != nil {
				return nil, err
			}

			stmt.AppendArg(val)
		}
		stmt.WriteString(binds)
	}
	if opt.Mode == options.InsertOnDuplicate {
		stmt.WriteString(" ON DUPLICATE KEY UPDATE ")
		next := false
		for _, f := range fields {
			if f.Path == pk {
				next = false
				continue
			}
			if next {
				stmt.WriteRune(',')
			}
			c := ms.Quote(f.Path)
			stmt.WriteString(c + "=VALUES(" + c + ")")
			next = true
		}
	}
	stmt.WriteRune(';')
	return
}

func findEncoder(mapper *reflext.Mapper, registry *codec.Registry, sf *reflext.StructField, v reflect.Value) (codec.ValueEncoder, error) {
	if _, ok := sf.Tag.LookUp("auto_increment"); ok && reflext.IsZero(v) {
		return codec.NilEncoder, nil
	}
	encoder, err := registry.LookupEncoder(v)
	if err != nil {
		return nil, err
	}
	return encoder, nil
}
