package mysql

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/internal/util"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/sql/charset"
	"github.com/si3nloong/sqlike/v2/sql/schema"
	sqltype "github.com/si3nloong/sqlike/v2/sql/type"
	sqlutil "github.com/si3nloong/sqlike/v2/sql/util"
	"github.com/si3nloong/sqlike/v2/x/reflext"
	"golang.org/x/text/currency"
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
	sb.SetTypeBuilder(sqltype.Time, s.TimeDataType)
	sb.SetTypeBuilder(sqltype.DateTime, s.DateTimeDataType)
	sb.SetTypeBuilder(sqltype.Timestamp, s.DateTimeDataType)
	sb.SetTypeBuilder(sqltype.UUID, s.UUIDDataType)
	sb.SetTypeBuilder(sqltype.JSON, s.JSONDataType)
	sb.SetTypeBuilder(sqltype.Point, s.SpatialDataType("POINT"))
	sb.SetTypeBuilder(sqltype.LineString, s.SpatialDataType("LINESTRING"))
	sb.SetTypeBuilder(sqltype.Polygon, s.SpatialDataType("POLYGON"))
	sb.SetTypeBuilder(sqltype.MultiPoint, s.SpatialDataType("MULTIPOINT"))
	sb.SetTypeBuilder(sqltype.MultiLineString, s.SpatialDataType("MULTILINESTRING"))
	sb.SetTypeBuilder(sqltype.MultiPolygon, s.SpatialDataType("MULTIPOLYGON"))
	sb.SetTypeBuilder(sqltype.String, s.StringDataType)
	sb.SetTypeBuilder(sqltype.Char, s.CharDataType)
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

func (s mySQLSchema) ByteDataType(sf reflext.FieldInfo) *sql.Column {
	col := new(sql.Column)
	col.Name = sf.Name()
	col.DataType = "MEDIUMBLOB"
	col.Type = "MEDIUMBLOB"
	col.Nullable = sf.IsNullable()
	tag := sf.Tag()
	if v, ok := tag.Option("default"); ok {
		col.DefaultValue = &v
	}
	return col
}

func (s mySQLSchema) UUIDDataType(sf reflext.FieldInfo) *sql.Column {
	col := new(sql.Column)
	col.Name = sf.Name()
	col.DataType = "BINARY"
	col.Type = "BINARY(16)"
	col.Size = 16
	col.Nullable = sf.IsNullable()
	return col
}

func (s mySQLSchema) DateDataType(sf reflext.FieldInfo) *sql.Column {
	col := new(sql.Column)
	col.Name = sf.Name()
	col.DataType = "DATE"
	col.Type = "DATE"
	col.Nullable = sf.IsNullable()
	return col
}

func (s mySQLSchema) TimeDataType(sf reflext.FieldInfo) *sql.Column {
	size := "6"
	if v, exists := sf.Tag().Option("size"); exists {
		if _, err := strconv.Atoi(v); err == nil {
			size = v
		}
	}

	col := new(sql.Column)
	col.Name = sf.Name()
	col.DataType = "TIME"
	col.Size = 6
	col.Type = "TIME(" + size + ")"
	col.Nullable = sf.IsNullable()
	// col.DefaultValue = &dflt
	// if _, ok := sf.Tag().Option("on_update"); ok {
	// 	col.Extra = "ON UPDATE " + dflt
	// }
	return col
}

func (s mySQLSchema) DateTimeDataType(sf reflext.FieldInfo) *sql.Column {
	size := "6"
	if v, exists := sf.Tag().Option("size"); exists {
		if _, err := strconv.Atoi(v); err == nil {
			size = v
		}
	}

	dflt := "CURRENT_TIMESTAMP(" + size + ")"
	col := new(sql.Column)
	col.Name = sf.Name()
	col.DataType = "DATETIME"
	col.Type = "DATETIME(" + size + ")"
	col.Nullable = sf.IsNullable()
	col.DefaultValue = &dflt
	if _, ok := sf.Tag().Option("on_update"); ok {
		col.Extra = "ON UPDATE " + dflt
	}
	return col
}

func (s mySQLSchema) JSONDataType(sf reflext.FieldInfo) *sql.Column {
	col := new(sql.Column)
	col.Name = sf.Name()
	col.DataType = "JSON"
	col.Type = "JSON"
	col.Nullable = sf.IsNullable()
	return col
}

func (s mySQLSchema) SpatialDataType(dataType string) schema.DataTypeFunc {
	return func(sf reflext.FieldInfo) *sql.Column {
		col := new(sql.Column)
		col.Name = sf.Name()
		col.DataType = dataType
		col.Type = dataType
		if sf.Type().Kind() == reflect.Ptr {
			col.Nullable = true
		}
		if v, ok := sf.Tag().Option("srid"); ok {
			if _, err := strconv.ParseUint(v, 10, 64); err != nil {
				panic("sqlike: SRID mosu be a number")
			}
			col.Extra = "SRID " + v
		}
		return col
	}
}

func (s mySQLSchema) StringDataType(sf reflext.FieldInfo) *sql.Column {
	col := new(sql.Column)
	col.Name = sf.Name()
	col.Nullable = sf.IsNullable()

	charset := "utf8mb4"
	collation := charsetMap[charset]
	dflt := ""
	tag := sf.Tag()
	cs, ok1 := tag.Option("charset")
	if ok1 {
		charset = strings.ToLower(cs)
		collation = charsetMap[charset]
	}

	col.DefaultValue = &dflt
	col.Charset = &charset
	col.Collation = &collation
	if v, ok := tag.Option("default"); ok {
		col.DefaultValue = &v
	}

	if enum, ok := tag.Option("enum"); ok {
		paths := strings.Split(enum, "|")
		if len(paths) < 1 {
			panic("invalid enum formats")
		}

		if !ok1 {
			charset = "utf8mb4"
			collation = "utf8mb4_unicode_ci"
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
		return col
	} else if char, ok := tag.Option("char"); ok {
		if _, err := strconv.Atoi(char); err != nil {
			panic("invalid value for char data type")
		}
		col.DataType = "CHAR"
		col.Type = "CHAR(" + char + ")"
		return col
	} else if _, ok := tag.Option("longtext"); ok {
		col.DataType = "TEXT"
		col.Type = "TEXT"
		col.DefaultValue = nil
		col.Charset = nil
		col.Collation = nil
		return col
	}

	size, _ := tag.Option("size")
	charLen, _ := strconv.Atoi(size)
	if charLen < 1 {
		charLen = 191
	}

	col.DataType = "VARCHAR"
	col.Type = "VARCHAR(" + strconv.Itoa(charLen) + ")"
	return col
}

func (s mySQLSchema) CharDataType(sf reflext.FieldInfo) *sql.Column {
	col := new(sql.Column)
	dflt := ""
	switch sf.Type() {
	case reflect.TypeOf(currency.Unit{}):
		charset, collation := string(charset.UTF8MB4), "utf8mb4_unicode_ci"
		col.Type = "CHAR(3)"
		col.Charset = &charset
		col.Collation = &collation
	default:
		charset, collation := string(charset.UTF8MB4), "utf8mb4_unicode_ci"
		col.Type = "CHAR(191)"
		col.Charset = &charset
		col.Collation = &collation
	}
	col.Name = sf.Name()
	col.DataType = "CHAR"
	col.Nullable = sf.IsNullable()
	col.DefaultValue = &dflt
	return col
}

func (s mySQLSchema) BoolDataType(sf reflext.FieldInfo) *sql.Column {
	col := new(sql.Column)
	dflt := "0"
	col.Name = sf.Name()
	col.DataType = "TINYINT"
	col.Type = "TINYINT(1)"
	col.Nullable = sf.IsNullable()
	col.DefaultValue = &dflt
	return col
}

func (s mySQLSchema) IntDataType(sf reflext.FieldInfo) *sql.Column {
	t := sf.Type()
	tag := sf.Tag()
	dflt := "0"
	dataType := s.getIntDataType(reflext.Deref(t))

	col := new(sql.Column)
	col.Name = sf.Name()
	col.DataType = dataType
	col.Type = dataType
	col.Nullable = sf.IsNullable()
	col.DefaultValue = &dflt
	if _, ok := tag.Option("auto_increment"); ok {
		col.Extra = "AUTO_INCREMENT"
		col.DefaultValue = nil
	} else if v, ok := tag.Option("default"); ok {
		if _, err := strconv.ParseUint(v, 10, 64); err != nil {
			panic("int default value should be integer")
		}
		col.DefaultValue = &v
	}
	return col
}

func (s mySQLSchema) UintDataType(sf reflext.FieldInfo) *sql.Column {
	t := sf.Type()
	tag := sf.Tag()
	dflt := "0"
	dataType := s.getIntDataType(reflext.Deref(t))

	col := new(sql.Column)
	col.Name = sf.Name()
	col.DataType = dataType
	col.Type = dataType + " UNSIGNED"
	col.Nullable = sf.IsNullable()
	col.DefaultValue = &dflt
	if _, ok := tag.Option("auto_increment"); ok {
		col.Extra = "AUTO_INCREMENT"
		col.DefaultValue = nil
	} else if v, ok := tag.Option("default"); ok {
		if _, err := strconv.ParseUint(v, 10, 64); err != nil {
			panic("uint default value should be unsigned integer")
		}
		col.DefaultValue = &v
	}
	return col
}

func (s mySQLSchema) FloatDataType(sf reflext.FieldInfo) *sql.Column {
	dflt := "0"
	tag := sf.Tag()

	col := new(sql.Column)
	col.Name = sf.Name()
	col.DataType = "REAL"
	col.Type = "REAL"
	if _, ok := tag.Option("unsigned"); ok {
		col.Type += " UNSIGNED"
	}
	col.Nullable = sf.IsNullable()
	col.DefaultValue = &dflt
	if v, ok := tag.Option("default"); ok {
		if _, err := strconv.ParseFloat(v, 64); err != nil {
			panic("float default value should be decimal number")
		}
		col.DefaultValue = &v
	}
	return col
}

func (s mySQLSchema) ArrayDataType(sf reflext.FieldInfo) *sql.Column {
	col := new(sql.Column)
	col.Name = sf.Name()
	col.Nullable = sf.IsNullable()
	// length := sf.Zero.Len()
	t := sf.Type().Elem()
	if t.Kind() == reflect.Uint8 {
		charset, collation := "ascii", "ascii_general_ci"
		col.DataType = "VARCHAR"
		col.Type = "VARCHAR(36)"
		col.Charset = &charset
		col.Collation = &collation
		return col
	}

	col.DataType = "JSON"
	col.Type = "JSON"
	return col
}

func (ms mySQL) buildSchemaByColumn(stmt db.Stmt, col *sql.Column) {
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
