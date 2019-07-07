package mysql

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/codec"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	sqlutil "github.com/si3nloong/sqlike/sql/util"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/primitive"
	"golang.org/x/xerrors"
)

var operatorMap = map[primitive.Operator]string{
	primitive.Equal:          "=",
	primitive.NotEqual:       "<>",
	primitive.Like:           "LIKE",
	primitive.NotLike:        "NOT LIKE",
	primitive.In:             "IN",
	primitive.NotIn:          "NOT IN",
	primitive.Between:        "BETWEEN",
	primitive.NotBetween:     "NOT BETWEEN",
	primitive.IsNull:         "IS NULL",
	primitive.NotNull:        "IS NOT NULL",
	primitive.GreaterThan:    ">",
	primitive.GreaterOrEqual: ">=",
	primitive.LesserThan:     "<",
	primitive.LesserOrEqual:  "<=",
	primitive.Or:             "OR",
	primitive.And:            "AND",
}

type mySQLBuilder struct {
	registry *codec.Registry
	builder  *sqlstmt.StatementBuilder
	sqlutil.MySQLUtil
}

func (b mySQLBuilder) SetRegistryAndBuilders(rg *codec.Registry, blr *sqlstmt.StatementBuilder) {
	if rg == nil {
		panic("missing required registry")
	}
	if blr == nil {
		panic("missing required parser")
	}
	blr.SetBuilder(reflect.TypeOf(primitive.L("")), b.ParseString)
	blr.SetBuilder(reflect.TypeOf(primitive.As{}), b.BuildAs)
	blr.SetBuilder(reflect.TypeOf(primitive.Raw{}), b.BuildRaw)
	blr.SetBuilder(reflect.TypeOf(primitive.Aggregate{}), b.BuildAggregate)
	blr.SetBuilder(reflect.TypeOf(primitive.Column{}), b.BuildColumn)
	blr.SetBuilder(reflect.TypeOf(primitive.C{}), b.ParseClause)
	blr.SetBuilder(reflect.TypeOf(primitive.Col("")), b.ParseString)
	blr.SetBuilder(reflect.TypeOf(primitive.Operator(0)), b.ParseOperator)
	blr.SetBuilder(reflect.TypeOf(primitive.G{}), b.BuildGroup)
	blr.SetBuilder(reflect.TypeOf(primitive.R{}), b.ParseRange)
	blr.SetBuilder(reflect.TypeOf(primitive.Sort{}), b.BuildSort)
	blr.SetBuilder(reflect.TypeOf(primitive.KV{}), b.ParseKeyValue)
	blr.SetBuilder(reflect.TypeOf(primitive.JC{}), b.ParseJSONContains)
	blr.SetBuilder(reflect.TypeOf(primitive.JQ("")), b.ParseJSONQuote)
	blr.SetBuilder(reflect.TypeOf(primitive.Math{}), b.ParseMath)
	blr.SetBuilder(reflect.TypeOf(&actions.FindActions{}), b.ParseFindActions)
	blr.SetBuilder(reflect.TypeOf(&actions.UpdateActions{}), b.ParseUpdateActions)
	blr.SetBuilder(reflect.TypeOf(&actions.DeleteActions{}), b.ParseDeleteActions)
	blr.SetBuilder(reflect.String, b.ParseString)
	b.registry = rg
	b.builder = blr
}

func (b *mySQLBuilder) ParseString(stmt *sqlstmt.Statement, it interface{}) error {
	v := reflect.ValueOf(it)
	stmt.WriteString(b.Quote(v.String()))
	return nil
}

// BuildColumn :
func (b *mySQLBuilder) BuildColumn(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.Column)
	if !isOk {
		return unmatchedDataType("BuildColumn")
	}
	stmt.WriteString(b.Quote(x.Name))
	return nil
}

func unmatchedDataType(callback string) error {
	return xerrors.New("invalid data type")
}

func (b *mySQLBuilder) BuildRaw(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.Raw)
	if !isOk {
		return unmatchedDataType("BuildRaw")
	}
	stmt.WriteString(x.Value)
	return nil
}

func (b *mySQLBuilder) BuildAs(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.As)
	if !isOk {
		return unmatchedDataType("BuildAs")
	}
	if err := b.getValue(stmt, x.Field); err != nil {
		return err
	}
	stmt.WriteString(" AS ")
	stmt.WriteString(b.Quote(x.Name))
	return nil
}

func (b *mySQLBuilder) BuildAggregate(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.Aggregate)
	if !isOk {
		return unmatchedDataType("BuildAggregate")
	}
	switch x.By {
	case primitive.Sum:
		stmt.WriteString("COALESCE(SUM(")
		if err := b.getValue(stmt, x.Field); err != nil {
			return err
		}
		stmt.WriteString("),0)")
		return nil
	case primitive.Average:
		stmt.WriteString("AVG")
	case primitive.Count:
		stmt.WriteString("COUNT")
	case primitive.Max:
		stmt.WriteString("MAX")
	case primitive.Min:
		stmt.WriteString("MIN")
	}
	stmt.WriteByte('(')
	if err := b.getValue(stmt, x.Field); err != nil {
		return err
	}
	stmt.WriteByte(')')
	return nil
}

func (b *mySQLBuilder) ParseOperator(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.Operator)
	if !isOk {
		return xerrors.New("operator")
	}
	stmt.WriteRune(' ')
	stmt.WriteString(operatorMap[x])
	stmt.WriteRune(' ')
	return nil
}

func (b *mySQLBuilder) ParseClause(stmt *sqlstmt.Statement, it interface{}) error {
	// stmt.WriteByte('(')
	x, isOk := it.(primitive.C)
	if !isOk {
		return unmatchedDataType("ParseClause")
	}

	if err := b.builder.BuildStatement(stmt, x.Field); err != nil {
		return err
	}

	stmt.WriteString(" " + operatorMap[x.Operator] + " ")
	switch x.Operator {
	case primitive.IsNull, primitive.NotNull:
		return nil
	}

	if err := b.getValue(stmt, x.Value); err != nil {
		return err
	}
	// stmt.WriteByte(')')

	// if x.Value == nil {
	// 	stmt.AppendArg(nil)
	// 	stmt.WriteRune('?')
	// 	return nil
	// }

	// v := reflext.ValueOf(x.Value)
	// if builder, isOk := b.builder.LookupBuilder(v.Type()); isOk {
	// 	if err := builder(stmt, x.Value); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }

	// encoder, err := b.registry.LookupEncoder(v.Type())
	// if err != nil {
	// 	return err
	// }
	// arg, err := encoder(nil, v)
	// if err != nil {
	// 	return err
	// }
	// stmt.AppendArg(arg)
	// stmt.WriteRune('?')
	return nil
}

func (b *mySQLBuilder) BuildSort(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.Sort)
	if !isOk {
		return unmatchedDataType("BuildSort")
	}
	stmt.WriteString(b.Quote(x.Field))
	if x.Order == primitive.Descending {
		stmt.WriteRune(' ')
		stmt.WriteString(`DESC`)
	}
	return nil
}

func (b *mySQLBuilder) ParseKeyValue(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x, isOk := it.(primitive.KV)
	if !isOk {
		return unmatchedDataType("ParseKeyValue")
	}
	stmt.WriteString(b.Quote(string(x.Field)))
	stmt.WriteString(` = `)

	return b.getValue(stmt, x.Value)
	// v := reflext.ValueOf(x.Value)
	// if !v.IsValid() {
	// 	stmt.WriteByte('?')
	// 	stmt.AppendArg(nil)
	// 	return
	// }

	// t := v.Type()
	// if builder, isOk := b.builder.LookupBuilder(t); isOk {
	// 	if err := builder(stmt, x.Value); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }

	// stmt.WriteByte('?')
	// encoder, err := b.registry.LookupEncoder(t)
	// if err != nil {
	// 	return err
	// }
	// vv, err := encoder(nil, v)
	// if err != nil {
	// 	return err
	// }
	// stmt.AppendArg(vv)
	// return
}

func (b *mySQLBuilder) ParseJSONContains(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.JC)
	if !isOk {
		return xerrors.New("expected json_contains")
	}
	stmt.WriteString("JSON_CONTAINS")
	stmt.WriteByte('(')
	args, err := b.encodeValue(x.Field)
	if err != nil {
		return err
	}
	stmt.WriteByte('?')
	stmt.AppendArg(args)
	stmt.WriteRune(',')
	if err := b.builder.BuildStatement(stmt, x.Value); err != nil {
		return err
	}
	stmt.WriteRune(',')
	stmt.WriteString("'" + x.Path + "'")
	stmt.WriteByte(')')
	return nil
}

func (b *mySQLBuilder) ParseJSONQuote(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.JQ)
	if !isOk {
		return xerrors.New("data type not match")
	}
	stmt.WriteString("JSON_QUOTE")
	stmt.WriteByte('(')
	stmt.WriteString(string(x))
	stmt.WriteByte(')')
	return nil
}

func (b *mySQLBuilder) ParseMath(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x, isOk := it.(primitive.Math)
	if !isOk {
		return unmatchedDataType("ParseMath")
	}
	// stmt.WriteByte('(')
	stmt.WriteString(b.Quote(string(x.Field)) + ` `)
	if x.Mode == primitive.Add {
		stmt.WriteRune('+')
	} else {
		stmt.WriteRune('-')
	}
	stmt.WriteString(` ` + strconv.Itoa(x.Value))
	// stmt.WriteByte(')')
	return
}

func (b *mySQLBuilder) ParseGroup(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x, isOk := it.(primitive.G)
	if !isOk {
		return xerrors.New("data type not match")
	}
	for len(x) > 0 {
		err = b.builder.BuildStatement(stmt, x[0])
		if err != nil {
			return
		}
		x = x[1:]
	}
	return
}

func (b *mySQLBuilder) getValue(stmt *sqlstmt.Statement, it interface{}) (err error) {
	v := reflext.ValueOf(it)
	if !v.IsValid() {
		stmt.WriteByte('?')
		stmt.AppendArg(nil)
		return
	}

	t := v.Type()
	if builder, isOk := b.builder.LookupBuilder(t); isOk {
		if err := builder(stmt, it); err != nil {
			return err
		}
		return nil
	}

	stmt.WriteByte('?')
	encoder, err := b.registry.LookupEncoder(t)
	if err != nil {
		return err
	}
	vv, err := encoder(nil, v)
	if err != nil {
		return err
	}
	stmt.AppendArg(vv)
	return
}

func (b *mySQLBuilder) BuildGroup(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x, isOk := it.(primitive.G)
	if !isOk {
		return unmatchedDataType("BuildGroup")
	}
	for len(x) > 0 {
		if err := b.getValue(stmt, x[0]); err != nil {
			return err
		}
		x = x[1:]
	}
	return
}

func (b *mySQLBuilder) encodeValue(it interface{}) (interface{}, error) {
	v := reflext.ValueOf(it)
	encoder, err := b.registry.LookupEncoder(v.Type())
	if err != nil {
		return nil, err
	}
	return encoder(nil, v)
}

func (b *mySQLBuilder) ParseRange(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x, isOk := it.(primitive.R)
	if !isOk {
		return xerrors.New("expected data type primitive.GV")
	}

	v := reflext.ValueOf(x.From)
	encoder, err := b.registry.LookupEncoder(v.Type())
	if err != nil {
		return err
	}
	arg, err := encoder(nil, v)
	if err != nil {
		return err
	}
	stmt.AppendArg(arg)

	v = reflext.ValueOf(x.To)
	encoder, err = b.registry.LookupEncoder(v.Type())
	if err != nil {
		return err
	}
	arg, err = encoder(nil, v)
	if err != nil {
		return err
	}
	stmt.AppendArg(arg)
	stmt.WriteRune('?')
	stmt.WriteString(" AND ")
	stmt.WriteRune('?')
	return
}

func (b *mySQLBuilder) ParseFindActions(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(*actions.FindActions)
	if !isOk {
		return xerrors.New("data type not match")
	}

	x.Table = strings.TrimSpace(x.Table)
	if x.Table == "" {
		return xerrors.New("empty table name")
	}
	stmt.WriteString(`SELECT `)
	if x.DistinctOn {
		stmt.WriteString("DISTINCT ")
	}
	if err := b.appendSelect(stmt, x.Projections); err != nil {
		return err
	}
	stmt.WriteString(` FROM ` + b.Quote(x.Table))
	if err := b.appendWhere(stmt, x.Conditions); err != nil {
		return err
	}
	if err := b.appendGroupBy(stmt, x.GroupBys); err != nil {
		return err
	}
	if err := b.appendOrderBy(stmt, x.Sorts); err != nil {
		return err
	}
	b.appendLimitNOffset(stmt, x.Record, x.Skip)

	return nil
}

func (b *mySQLBuilder) ParseUpdateActions(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(*actions.UpdateActions)
	if !isOk {
		return xerrors.New("data type not match")
	}
	stmt.WriteString(`UPDATE ` + b.Quote(x.Table) + ` `)
	if err := b.appendSet(stmt, x.Values); err != nil {
		return err
	}
	if err := b.appendWhere(stmt, x.Conditions); err != nil {
		return err
	}
	if err := b.appendOrderBy(stmt, x.Sorts); err != nil {
		return err
	}
	b.appendLimitNOffset(stmt, x.Record, 0)
	return nil
}

func (b *mySQLBuilder) ParseDeleteActions(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(*actions.DeleteActions)
	if !isOk {
		return xerrors.New("data type not match")
	}
	stmt.WriteString(`DELETE FROM ` + b.Quote(x.Table))
	if err := b.appendWhere(stmt, x.Conditions); err != nil {
		return err
	}
	if err := b.appendOrderBy(stmt, x.Sorts); err != nil {
		return err
	}
	b.appendLimitNOffset(stmt, x.Record, 0)
	return nil
}

func (b *mySQLBuilder) appendSelect(stmt *sqlstmt.Statement, pjs []interface{}) error {
	if len(pjs) > 0 {
		length := len(pjs)
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteRune(',')
			}
			if err := b.builder.BuildStatement(stmt, pjs[i]); err != nil {
				return err
			}
		}
		return nil
	}

	stmt.WriteString(`*`)
	return nil
}

func (b *mySQLBuilder) appendWhere(stmt *sqlstmt.Statement, conds []interface{}) error {
	length := len(conds)
	if length > 0 {
		stmt.WriteString(` WHERE `)
		for i := 0; i < length; i++ {
			if err := b.builder.BuildStatement(stmt, conds[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *mySQLBuilder) appendGroupBy(stmt *sqlstmt.Statement, fields []interface{}) error {
	length := len(fields)
	if length > 0 {
		stmt.WriteString(` GROUP BY `)
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteRune(',')
			}
			if err := b.builder.BuildStatement(stmt, fields[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *mySQLBuilder) appendOrderBy(stmt *sqlstmt.Statement, sorts []primitive.Sort) error {
	length := len(sorts)
	if length > 0 {
		stmt.WriteString(` ORDER BY `)
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteRune(',')
			}
			if err := b.builder.BuildStatement(stmt, sorts[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *mySQLBuilder) appendLimitNOffset(stmt *sqlstmt.Statement, limit, offset uint) {
	if limit > 0 {
		stmt.WriteString(` LIMIT ` + strconv.FormatUint(uint64(limit), 10))
	}
	if offset > 0 {
		stmt.WriteString(` OFFSET ` + strconv.FormatUint(uint64(offset), 10))
	}
}

func (b *mySQLBuilder) appendSet(stmt *sqlstmt.Statement, values []primitive.KV) error {
	length := len(values)
	if length > 0 {
		stmt.WriteString(`SET `)
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteRune(',')
			}
			if err := b.builder.BuildStatement(stmt, values[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
