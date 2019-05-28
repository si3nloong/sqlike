package mysql

import (
	"reflect"
	"strconv"
	"strings"

	"bitbucket.org/SianLoong/sqlike/reflext"
	"bitbucket.org/SianLoong/sqlike/sqlike/sql/component"
	"bitbucket.org/SianLoong/sqlike/sqlike/sql/internal"
	sqltype "bitbucket.org/SianLoong/sqlike/sqlike/sql/types"
	sqlutil "bitbucket.org/SianLoong/sqlike/sqlike/sql/util"
	"bitbucket.org/SianLoong/sqlike/util"
	"golang.org/x/xerrors"
)

// mySQLSchema :
type mySQLSchema struct {
	sqlutil.MySQLUtil
}

// SetBuilders :
func (s mySQLSchema) SetBuilders(sb *internal.SchemaBuilder) {
	sb.SetTypeBuilder(sqltype.Byte, s.ByteDataType)
	sb.SetTypeBuilder(sqltype.Date, s.DateDataType)
	sb.SetTypeBuilder(sqltype.DateTime, s.TimeDataType)
	sb.SetTypeBuilder(sqltype.Timestamp, s.TimeDataType)
	sb.SetTypeBuilder(sqltype.JSON, s.JSONDataType)
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
	sb.SetTypeBuilder(sqltype.Array, s.JSONDataType)
	sb.SetTypeBuilder(sqltype.Slice, s.JSONDataType)
	sb.SetTypeBuilder(sqltype.Map, s.JSONDataType)
}

func (s mySQLSchema) ByteDataType(sf *reflext.StructField) (col component.Column) {
	col.Name = sf.Path
	col.DataType = `MEDIUMBLOB`
	col.Type = `MEDIUMBLOB`
	col.Nullable = sf.IsNullable
	return
}

func (s mySQLSchema) DateDataType(sf *reflext.StructField) (col component.Column) {
	dflt := `CURDATE()`
	col.Name = sf.Path
	col.DataType = `DATE`
	col.Type = `DATE`
	col.Nullable = sf.IsNullable
	col.DefaultValue = &dflt
	return
}

func (s mySQLSchema) TimeDataType(sf *reflext.StructField) (col component.Column) {
	dflt := `CURRENT_TIMESTAMP(6)`
	col.Name = sf.Path
	col.DataType = `DATETIME(6)`
	col.Type = `DATETIME(6)`
	col.Nullable = sf.IsNullable
	col.DefaultValue = &dflt
	return
}

func (s mySQLSchema) JSONDataType(sf *reflext.StructField) (col component.Column) {
	col.Name = sf.Path
	col.DataType = `JSON`
	col.Type = `JSON`
	col.Nullable = sf.IsNullable
	return
}

func (s mySQLSchema) StringDataType(sf *reflext.StructField) (col component.Column) {
	col.Name = sf.Path
	col.Nullable = sf.IsNullable

	if enum, isOk := sf.Tag.LookUp("enum"); isOk {
		paths := strings.Split(enum, "|")
		if len(paths) < 1 {
			panic(xerrors.New("invalid enum formats"))
		}

		blr := util.AcquireString()
		defer util.ReleaseString(blr)
		blr.WriteString(`ENUM`)
		blr.WriteRune('(')
		for i, p := range paths {
			if i > 0 {
				blr.WriteRune(',')
			}
			blr.WriteString(s.Wrap(p))
		}
		blr.WriteRune(')')

		dflt := paths[0]
		col.DataType = `ENUM`
		col.Type = blr.String()
		col.DefaultValue = &dflt
		return
	} else if _, isOk := sf.Tag.LookUp("longtext"); isOk {
		col.DataType = `TEXT`
		col.Type = `TEXT`
		return
	}

	size, _ := sf.Tag.LookUp("size")
	charLen, _ := strconv.Atoi(size)
	if charLen < 1 {
		charLen = 191
	}

	dflt := ``
	charset := `utf8mb4`
	collation := `utf8mb4_unicode_ci`

	col.DefaultValue = &dflt
	col.DataType = `VARCHAR`
	col.Type = `VARCHAR(` + strconv.Itoa(charLen) + `)`
	col.CharSet = &charset
	col.Collation = &collation
	return
}

func (s mySQLSchema) BoolDataType(sf *reflext.StructField) (col component.Column) {
	dflt := `0`
	col.Name = sf.Path
	col.DataType = `TINYINT`
	col.Type = `TINYINT(1)`
	col.Nullable = sf.IsNullable
	col.DefaultValue = &dflt
	return
}

func (s mySQLSchema) IntDataType(sf *reflext.StructField) (col component.Column) {
	t := sf.Zero.Type()
	dflt := `0`
	dataType := s.getIntDataType(reflext.Deref(t))

	col.Name = sf.Path
	col.DataType = dataType
	col.Type = dataType
	col.Nullable = sf.IsNullable
	if _, isOk := sf.Tag.LookUp("auto_increment"); isOk {
		col.Extra = "AUTO_INCREMENT"
	} else {
		col.DefaultValue = &dflt
	}
	return
}

func (s mySQLSchema) UintDataType(sf *reflext.StructField) (col component.Column) {
	t := sf.Zero.Type()
	dflt := `0`
	dataType := s.getIntDataType(reflext.Deref(t))

	col.Name = sf.Path
	col.DataType = dataType
	col.Type = dataType + ` UNSIGNED`
	col.Nullable = sf.IsNullable
	if _, isOk := sf.Tag.LookUp("auto_increment"); isOk {
		col.Extra = "AUTO_INCREMENT"
	} else {
		col.DefaultValue = &dflt
	}
	return
}

func (s mySQLSchema) FloatDataType(sf *reflext.StructField) (col component.Column) {
	dflt := `0`
	col.Name = sf.Path
	col.DataType = `REAL`
	col.Type = `REAL`
	if _, isOk := sf.Tag.LookUp("unsigned"); isOk {
		col.Type += ` UNSIGNED`
	}
	col.Nullable = sf.IsNullable
	col.DefaultValue = &dflt
	return
}

func (s mySQLSchema) getIntDataType(t reflect.Type) (dataType string) {
	switch t.Kind() {
	case reflect.Int8, reflect.Uint8:
		dataType = `TINYINT`
	case reflect.Int16, reflect.Uint16:
		dataType = `SMALLINT`
	case reflect.Int32, reflect.Uint32:
		dataType = `INT`
	case reflect.Int64, reflect.Uint64:
		dataType = `BIGINT`
	default:
		dataType = `INT`
	}
	return
}
