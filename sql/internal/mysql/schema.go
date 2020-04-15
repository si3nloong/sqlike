package mysql

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/schema"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	sqltype "github.com/si3nloong/sqlike/sql/type"
	sqlutil "github.com/si3nloong/sqlike/sql/util"
	"github.com/si3nloong/sqlike/sqlike/columns"
	"github.com/si3nloong/sqlike/util"
)

var charsetMap = map[string]string{
	"utf8mb4": "utf8mb4_unicode_ci",
	"latin1":  "latin1_bin",
}

// mySQLSchema :
type mySQLSchema struct {
	sqlutil.MySQLUtil
}

// SetBuilders :
func (s mySQLSchema) SetBuilders(sb *schema.Builder) {
	sb.SetTypeBuilder(sqltype.Byte, s.ByteDataType)
	sb.SetTypeBuilder(sqltype.Date, s.DateDataType)
	sb.SetTypeBuilder(sqltype.DateTime, s.TimeDataType)
	sb.SetTypeBuilder(sqltype.Timestamp, s.TimeDataType)
	sb.SetTypeBuilder(sqltype.UUID, s.UUIDDataType)
	sb.SetTypeBuilder(sqltype.JSON, s.JSONDataType)
	sb.SetTypeBuilder(sqltype.Point, s.SpatialDataType("POINT"))
	sb.SetTypeBuilder(sqltype.LineString, s.SpatialDataType("LINESTRING"))
	sb.SetTypeBuilder(sqltype.Polygon, s.SpatialDataType("POLYGON"))
	sb.SetTypeBuilder(sqltype.MultiPoint, s.SpatialDataType("MULTIPOINT"))
	sb.SetTypeBuilder(sqltype.MultiLineString, s.SpatialDataType("MULTILINESTRING"))
	sb.SetTypeBuilder(sqltype.MultiPolygon, s.SpatialDataType("MULTIPOLYGON"))
	sb.SetTypeBuilder(sqltype.String, s.StringDataType)
	sb.SetTypeBuilder(sqltype.Bool, s.BoolDataType)
	sb.SetTypeBuilder(sqltype.Int, s.IntDataType)
	sb.SetTypeBuilder(sqltype.Int8, s.IntDataType)
	sb.SetTypeBuilder(sqltype.Int16, s.IntDataType)
	sb.SetTypeBuilder(sqltype.Int32, s.IntDataType)
	sb.SetTypeBuilder(sqltype.Int64, s.IntDataType)
	sb.SetTypeBuilder(sqltype.Uint, s.UintDataType)
	sb.SetTypeBuilder(sqltype.Uint8, s.UintDataType)
	sb.SetTypeBuilder(sqltype.Uint16, s.UintDataType)
	sb.SetTypeBuilder(sqltype.Uint32, s.UintDataType)
	sb.SetTypeBuilder(sqltype.Uint64, s.UintDataType)
	sb.SetTypeBuilder(sqltype.Float32, s.FloatDataType)
	sb.SetTypeBuilder(sqltype.Float64, s.FloatDataType)
	sb.SetTypeBuilder(sqltype.Struct, s.JSONDataType)
	sb.SetTypeBuilder(sqltype.Array, s.ArrayDataType)
	sb.SetTypeBuilder(sqltype.Slice, s.JSONDataType)
	sb.SetTypeBuilder(sqltype.Map, s.JSONDataType)
}

func (s mySQLSchema) ByteDataType(sf *reflext.StructField) (col columns.Column) {
	col.Name = sf.Path
	col.DataType = "MEDIUMBLOB"
	col.Type = "MEDIUMBLOB"
	col.Nullable = sf.IsNullable
	return
}

func (s mySQLSchema) UUIDDataType(sf *reflext.StructField) (col columns.Column) {
	charset, collation := "ascii", "ascii_general_ci"
	col.Name = sf.Path
	col.DataType = "VARCHAR"
	col.Type = "VARCHAR(36)"
	col.Size = 36
	col.Charset = &charset
	col.Collation = &collation
	col.Nullable = sf.IsNullable
	return
}

func (s mySQLSchema) DateDataType(sf *reflext.StructField) (col columns.Column) {
	dflt := "CURDATE()"
	col.Name = sf.Path
	col.DataType = "DATE"
	col.Type = "DATE"
	col.Nullable = sf.IsNullable
	col.DefaultValue = &dflt
	return
}

func (s mySQLSchema) TimeDataType(sf *reflext.StructField) (col columns.Column) {
	size := "6"
	if v, exists := sf.Tag.LookUp("size"); exists {
		if _, err := strconv.Atoi(v); err == nil {
			size = v
		}
	}

	dflt := "CURRENT_TIMESTAMP(" + size + ")"
	col.Name = sf.Path
	col.DataType = "DATETIME"
	col.Type = "DATETIME(" + size + ")"
	col.Nullable = sf.IsNullable
	col.DefaultValue = &dflt
	return
}

func (s mySQLSchema) JSONDataType(sf *reflext.StructField) (col columns.Column) {
	col.Name = sf.Path
	col.DataType = "JSON"
	col.Type = "JSON"
	col.Nullable = sf.IsNullable
	return
}

func (s mySQLSchema) SpatialDataType(dataType string) schema.DataTypeFunc {
	return func(sf *reflext.StructField) (col columns.Column) {
		col.Name = sf.Path
		col.DataType = dataType
		col.Type = dataType
		if sf.Type.Kind() == reflect.Ptr {
			col.Nullable = true
		}
		if v, ok := sf.Tag.LookUp("srid"); ok {
			if _, err := strconv.ParseUint(v, 10, 64); err != nil {
				return
			}
			col.Extra = "SRID " + v
		}
		return
	}
}

func (s mySQLSchema) StringDataType(sf *reflext.StructField) (col columns.Column) {
	col.Name = sf.Path
	col.Nullable = sf.IsNullable

	charset := "utf8mb4"
	collation := charsetMap[charset]
	dflt := ""
	cs, ok1 := sf.Tag.LookUp("charset")
	if ok1 {
		charset = strings.ToLower(cs)
		collation = charsetMap[charset]
	}

	col.DefaultValue = &dflt
	col.Charset = &charset
	col.Collation = &collation

	if enum, ok := sf.Tag.LookUp("enum"); ok {
		paths := strings.Split(enum, "|")
		if len(paths) < 1 {
			panic("invalid enum formats")
		}

		// if charset is not set, default is using latin1 for better performance
		if !ok1 {
			charset = "latin1"
			collation = "latin1_bin"
		}

		blr := util.AcquireString()
		defer util.ReleaseString(blr)
		blr.WriteString("ENUM")
		blr.WriteRune('(')
		for i, p := range paths {
			if i > 0 {
				blr.WriteRune(',')
			}
			blr.WriteString(s.Wrap(p))
		}
		blr.WriteRune(')')

		dflt = paths[0]
		col.DataType = "ENUM"
		col.Type = blr.String()
		col.DefaultValue = &dflt
		return
	} else if char, ok := sf.Tag.LookUp("char"); ok {
		if _, err := strconv.Atoi(char); err != nil {
			panic("invalid value for char data type")
		}
		col.DataType = "CHAR(" + char + ")"
		col.Type = "CHAR(" + char + ")"
		return
	} else if _, ok := sf.Tag.LookUp("longtext"); ok {
		col.DataType = "TEXT"
		col.Type = "TEXT"
		col.DefaultValue = nil
		col.Charset = nil
		col.Collation = nil
		return
	}

	size, _ := sf.Tag.LookUp("size")
	charLen, _ := strconv.Atoi(size)
	if charLen < 1 {
		charLen = 191
	}

	col.DataType = "VARCHAR"
	col.Type = "VARCHAR(" + strconv.Itoa(charLen) + ")"
	return
}

func (s mySQLSchema) BoolDataType(sf *reflext.StructField) (col columns.Column) {
	dflt := "0"
	col.Name = sf.Path
	col.DataType = "TINYINT"
	col.Type = "TINYINT(1)"
	col.Nullable = sf.IsNullable
	col.DefaultValue = &dflt
	return
}

func (s mySQLSchema) IntDataType(sf *reflext.StructField) (col columns.Column) {
	t := sf.Type
	dflt := "0"
	dataType := s.getIntDataType(reflext.Deref(t))

	col.Name = sf.Path
	col.DataType = dataType
	col.Type = dataType
	col.Nullable = sf.IsNullable
	col.DefaultValue = &dflt
	if _, ok := sf.Tag.LookUp("auto_increment"); ok {
		col.Extra = "AUTO_INCREMENT"
		col.DefaultValue = nil
	}
	return
}

func (s mySQLSchema) UintDataType(sf *reflext.StructField) (col columns.Column) {
	t := sf.Type
	dflt := "0"
	dataType := s.getIntDataType(reflext.Deref(t))

	col.Name = sf.Path
	col.DataType = dataType
	col.Type = dataType + " UNSIGNED"
	col.Nullable = sf.IsNullable
	col.DefaultValue = &dflt
	if _, ok := sf.Tag.LookUp("auto_increment"); ok {
		col.Extra = "AUTO_INCREMENT"
		col.DefaultValue = nil
	}
	return
}

func (s mySQLSchema) FloatDataType(sf *reflext.StructField) (col columns.Column) {
	dflt := "0"
	col.Name = sf.Path
	col.DataType = "REAL"
	col.Type = "REAL"
	if _, ok := sf.Tag.LookUp("unsigned"); ok {
		col.Type += " UNSIGNED"
	}
	col.Nullable = sf.IsNullable
	col.DefaultValue = &dflt
	return
}

func (s mySQLSchema) ArrayDataType(sf *reflext.StructField) (col columns.Column) {
	col.Name = sf.Path
	col.Nullable = sf.IsNullable
	// length := sf.Zero.Len()
	t := sf.Type.Elem()
	if t.Kind() == reflect.Uint8 {
		charset, collation := "ascii", "ascii_general_ci"
		col.DataType = "VARCHAR"
		col.Type = "VARCHAR(36)"
		col.Charset = &charset
		col.Collation = &collation
		return
	}
	col.DataType = "JSON"
	col.Type = "JSON"
	return
}

func (ms MySQL) buildSchemaByColumn(stmt *sqlstmt.Statement, col columns.Column) {
	stmt.WriteString(ms.Quote(col.Name))
	stmt.WriteString(" " + col.Type)
	if col.Charset != nil {
		stmt.WriteString(" CHARACTER SET " + *col.Charset)
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

func (s mySQLSchema) getIntDataType(t reflect.Type) (dataType string) {
	switch t.Kind() {
	case reflect.Int8, reflect.Uint8:
		dataType = "TINYINT"
	case reflect.Int16, reflect.Uint16:
		dataType = "SMALLINT"
	case reflect.Int32, reflect.Uint32:
		dataType = "INT"
	case reflect.Int64, reflect.Uint64:
		dataType = "BIGINT"
	default:
		dataType = "INT"
	}
	return
}
