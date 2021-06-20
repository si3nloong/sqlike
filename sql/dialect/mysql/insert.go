package mysql

import (
	"fmt"
	"reflect"

	"github.com/si3nloong/sqlike/db"
	"github.com/si3nloong/sqlike/options"
	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/x/reflext"
	"github.com/si3nloong/sqlike/x/spatial"
)

// InsertInto :
func (ms MySQL) InsertInto(stmt db.Stmt, db, table, pk string, cache reflext.StructMapper, cdc codec.Codecer, fields []reflext.StructFielder, v reflect.Value, opt *options.InsertOptions) (err error) {
	records := v.Len()

	stmt.WriteString("INSERT")
	if opt.Mode == options.InsertIgnore {
		stmt.WriteString(" IGNORE")
	}
	stmt.WriteString(" INTO " + ms.TableName(db, table) + " (")

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
		if _, ok := fields[i].Tag().LookUp("generated_column"); ok {
			fields = append(fields[:i], fields[i+1:]...)
			continue
		}

		stmt.WriteString(ms.Quote(fields[i].Name()))
		if i < len(fields)-1 {
			stmt.WriteByte(',')
		}

		i++
	}
	stmt.WriteString(") VALUES ")

	length := len(fields)
	encoders := make([]codec.ValueEncoder, length)
	for i := 0; i < records; i++ {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteByte('(')
		vi := reflext.Indirect(v.Index(i))

		for j := range fields {
			if j > 0 {
				stmt.WriteByte(',')
			}

			// first record only find encoders
			fv := cache.FieldByIndexesReadOnly(vi, fields[j].Index())
			if i == 0 {
				encoders[j], err = findEncoder(cdc, fields[j], fv)
				if err != nil {
					return err
				}
			}

			val, err := encoders[j](fields[j], fv)
			if err != nil {
				return err
			}

			convertSpatial(stmt, val)
		}
		stmt.WriteByte(')')
	}

	var column string
	if opt.Mode == options.InsertOnDuplicate {
		stmt.WriteString(" ON DUPLICATE KEY UPDATE ")
		next := false
		for i := range fields {
			// skip primary key on duplicate update
			if fields[i].Name() == pk {
				next = false
				continue
			}

			// skip omit fields on update
			if _, ok := omitField[fields[i].Name()]; ok {
				continue
			}
			if next {
				stmt.WriteByte(',')
			}

			column = ms.Quote(fields[i].Name())
			stmt.WriteString(column + "=VALUES(" + column + ")")
			next = true
		}
	}
	stmt.WriteByte(';')
	return
}

func findEncoder(c codec.Codecer, sf reflext.StructFielder, v reflect.Value) (codec.ValueEncoder, error) {
	// auto_increment field should pass nil if it's empty
	if _, ok := sf.Tag().LookUp("auto_increment"); ok && reflext.IsZero(v) {
		return codec.NilEncoder, nil
	}
	encoder, err := c.LookupEncoder(v)
	if err != nil {
		return nil, err
	}
	return encoder, nil
}

func convertSpatial(stmt db.Stmt, val interface{}) {
	switch vi := val.(type) {
	case spatial.Geometry:
		switch vi.Type {
		case spatial.Point:
			stmt.WriteString("ST_PointFromText")
		case spatial.LineString:
			stmt.WriteString("ST_LineStringFromText")
		case spatial.Polygon:
			stmt.WriteString("ST_PolygonFromText")
		case spatial.MultiPoint:
			stmt.WriteString("ST_MultiPointFromText")
		case spatial.MultiLineString:
			stmt.WriteString("ST_MultiLineStringFromText")
		case spatial.MultiPolygon:
			stmt.WriteString("ST_MultiPolygonFromText")
		default:
		}

		stmt.WriteString("(?")
		if vi.SRID > 0 {
			stmt.WriteString(fmt.Sprintf(",%d", vi.SRID))
		}
		stmt.WriteByte(')')
		stmt.AppendArgs(vi.WKT)

	default:
		stmt.WriteByte('?')
		stmt.AppendArgs(val)
	}
}
