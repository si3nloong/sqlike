package mysql

import (
	"fmt"
	"reflect"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/spatial"
	"github.com/si3nloong/sqlike/sql/codec"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// InsertInto :
func (ms MySQL) InsertInto(stmt sqlstmt.Stmt, db, table, pk string, cache reflext.StructMapper, cdc codec.Codecer, fields []reflext.StructFielder, v reflect.Value, opt *options.InsertOptions) (err error) {
	records := v.Len()
	stmt.WriteString("INSERT")
	if opt.Mode == options.InsertIgnore {
		stmt.WriteString(" IGNORE")
	}
	stmt.WriteString(" INTO " + ms.TableName(db, table) + " (")
	for i, f := range fields {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(ms.Quote(f.Name()))
	}
	stmt.WriteString(") VALUES ")
	length := len(fields)
	// binds := strings.Repeat("?,", length)
	// binds = "(" + binds[:len(binds)-1] + ")"
	encoders := make([]codec.ValueEncoder, length)
	for i := 0; i < records; i++ {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteByte('(')
		vi := reflext.Indirect(v.Index(i))
		for j, sf := range fields {
			if j > 0 {
				stmt.WriteByte(',')
			}
			// first record only find encoders
			fv := cache.FieldByIndexesReadOnly(vi, sf.Index())
			if i == 0 {
				encoders[j], err = findEncoder(cdc, sf, fv)
				if err != nil {
					return err
				}
			}

			val, err := encoders[j](sf, fv)
			if err != nil {
				return err
			}

			convertSpatial(stmt, val)

		}
		stmt.WriteByte(')')
		// stmt.WriteString(binds)
	}
	if opt.Mode == options.InsertOnDuplicate {
		stmt.WriteString(" ON DUPLICATE KEY UPDATE ")
		next := false
		for _, f := range fields {
			if f.Name() == pk {
				next = false
				continue
			}
			if next {
				stmt.WriteByte(',')
			}
			c := ms.Quote(f.Name())
			stmt.WriteString(c + "=VALUES(" + c + ")")
			next = true
		}
	}
	stmt.WriteByte(';')
	return
}

func findEncoder(c codec.Codecer, sf reflext.StructFielder, v reflect.Value) (codec.ValueEncoder, error) {
	if _, ok := sf.Tag().LookUp("auto_increment"); ok && reflext.IsZero(v) {
		return codec.NilEncoder, nil
	}
	encoder, err := c.LookupEncoder(v)
	if err != nil {
		return nil, err
	}
	return encoder, nil
}

func convertSpatial(stmt sqlstmt.Stmt, val interface{}) {
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
