package mysql

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/codec"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	sqlutil "github.com/si3nloong/sqlike/sql/util"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

var operatorMap = map[primitive.Operator]string{
	primitive.Equal:          "=",
	primitive.NotEqual:       "<>",
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
	blr.SetBuilder(reflect.TypeOf(primitive.CastAs{}), b.BuildCastAs)
	blr.SetBuilder(reflect.TypeOf(primitive.Func{}), b.BuildFunction)
	blr.SetBuilder(reflect.TypeOf(primitive.JSONFunc{}), b.BuildJSONFunction)
	blr.SetBuilder(reflect.TypeOf(primitive.Field{}), b.BuildField)
	blr.SetBuilder(reflect.TypeOf(primitive.Value{}), b.BuildValue)
	blr.SetBuilder(reflect.TypeOf(primitive.As{}), b.BuildAs)
	blr.SetBuilder(reflect.TypeOf(primitive.Nil{}), b.BuildNil)
	blr.SetBuilder(reflect.TypeOf(primitive.Raw{}), b.BuildRaw)
	blr.SetBuilder(reflect.TypeOf(primitive.Aggregate{}), b.BuildAggregate)
	blr.SetBuilder(reflect.TypeOf(primitive.Column{}), b.BuildColumn)
	blr.SetBuilder(reflect.TypeOf(primitive.JSONColumn{}), b.BuildJSONColumn)
	blr.SetBuilder(reflect.TypeOf(primitive.C{}), b.BuildClause)
	blr.SetBuilder(reflect.TypeOf(primitive.L{}), b.BuildLike)
	blr.SetBuilder(reflect.TypeOf(primitive.Operator(0)), b.BuildOperator)
	blr.SetBuilder(reflect.TypeOf(primitive.Group{}), b.BuildGroup)
	blr.SetBuilder(reflect.TypeOf(primitive.R{}), b.BuildRange)
	blr.SetBuilder(reflect.TypeOf(primitive.Sort{}), b.BuildSort)
	blr.SetBuilder(reflect.TypeOf(primitive.KV{}), b.BuildKeyValue)
	blr.SetBuilder(reflect.TypeOf(primitive.JC{}), b.BuildJSONContains)
	blr.SetBuilder(reflect.TypeOf(primitive.Math{}), b.BuildMath)
	blr.SetBuilder(reflect.TypeOf(&actions.FindActions{}), b.BuildFindActions)
	blr.SetBuilder(reflect.TypeOf(&actions.UpdateActions{}), b.BuildUpdateActions)
	blr.SetBuilder(reflect.TypeOf(&actions.DeleteActions{}), b.BuildDeleteActions)
	blr.SetBuilder(reflect.String, b.BuildString)
	b.registry = rg
	b.builder = blr
}

func (b *mySQLBuilder) BuildCastAs(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.CastAs)
	stmt.WriteString("CAST(")
	if err := b.builder.BuildStatement(stmt, x.Value); err != nil {
		return err
	}
	stmt.WriteString(" AS ")
	switch x.DataType {
	case primitive.JSON:
		stmt.WriteString("JSON")
	default:
		return errors.New("mysql: unsupported cast as data type")
	}
	stmt.WriteByte(')')
	return nil
}

func (b *mySQLBuilder) BuildFunction(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.Func)
	stmt.WriteString(x.Name)
	stmt.WriteByte('(')
	for i, args := range x.Arguments {
		if i > 0 {
			stmt.WriteByte(',')
		}
		if err := b.builder.BuildStatement(stmt, args); err != nil {
			return err
		}
	}
	stmt.WriteByte(')')
	return nil
}

func (b *mySQLBuilder) BuildJSONFunction(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.JSONFunc)
	switch x.Type {
	case primitive.JSONQuote:
		stmt.WriteString("JSON_QUOTE")
	case primitive.JSONContains:

	default:
		return errors.New("mysql: unsupported function")
	}
	stmt.WriteByte('(')
	for i, args := range x.Arguments {
		if i > 0 {
			stmt.WriteByte(',')
		}
		if err := b.builder.BuildStatement(stmt, args); err != nil {
			return err
		}
	}
	stmt.WriteByte(')')
	return nil
}

func (b *mySQLBuilder) BuildString(stmt *sqlstmt.Statement, it interface{}) error {
	v := reflect.ValueOf(it)
	stmt.WriteString(b.Quote(v.String()))
	return nil
}

func (b *mySQLBuilder) BuildLike(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.L)
	if err := b.builder.BuildStatement(stmt, x.Field); err != nil {
		return err
	}
	stmt.WriteByte(' ')
	if x.IsNot {
		stmt.WriteString("NOT LIKE")
	} else {
		stmt.WriteString("LIKE")
	}
	stmt.WriteByte(' ')
	v := reflext.ValueOf(x.Value)
	if !v.IsValid() {
		stmt.WriteByte('?')
		stmt.AppendArg(nil)
		return nil
	}

	t := v.Type()
	if builder, ok := b.builder.LookupBuilder(t); ok {
		if err := builder(stmt, it); err != nil {
			return err
		}
		return nil
	}

	stmt.WriteByte('?')
	encoder, err := b.registry.LookupEncoder(v)
	if err != nil {
		return err
	}
	vv, err := encoder(nil, v)
	if err != nil {
		return err
	}
	switch vi := vv.(type) {
	case string:
		vv = escapeWildCard(vi)
	case []byte:
		vv = escapeWildCard(string(vi))
	}
	stmt.AppendArg(vv)
	return nil
}

func (b *mySQLBuilder) BuildField(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.Field)
	stmt.WriteString("FIELD")
	stmt.WriteByte('(')
	stmt.WriteString(b.Quote(x.Name))
	for _, v := range x.Values {
		stmt.WriteByte(',')
		if err := b.getValue(stmt, v); err != nil {
			return err
		}
	}
	stmt.WriteByte(')')
	return nil
}

func (b *mySQLBuilder) BuildValue(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x := it.(primitive.Value)
	v := reflext.ValueOf(x.Raw)
	if !v.IsValid() {
		stmt.WriteByte('?')
		stmt.AppendArg(nil)
		return
	}

	stmt.WriteByte('?')
	encoder, err := b.registry.LookupEncoder(v)
	if err != nil {
		return err
	}
	vv, err := encoder(nil, v)
	if err != nil {
		return err
	}
	stmt.AppendArg(vv)
	return nil
}

// BuildColumn :
func (b *mySQLBuilder) BuildColumn(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.Column)
	stmt.WriteString(b.Quote(x.Name))
	return nil
}

// BuildJSONColumn :
func (b *mySQLBuilder) BuildJSONColumn(stmt *sqlstmt.Statement, it interface{}) error {
	/*
		Expected columns ( JSON_EXTRACT )
		Column : Address
		Nested : [ State, City ]
		UnquoteResult : false

		Result
		`Address`->'$.State.City'

		--------------------------------------------

		Expected columns ( JSON_EXTRACT(JSON_UNQUOTE) )
		Column : Address
		Nested : [ State, City ]
		UnquoteResult : true

		Result
		`Address`->>'$.State.City'
	*/
	x := it.(primitive.JSONColumn)
	nested := strings.Join(x.Nested, ".")
	operator := "->"
	if !strings.HasPrefix(nested, "$.") {
		nested = "$." + nested
	}
	if x.UnquoteResult {
		operator += ">"
	}
	stmt.WriteString(b.Quote(x.Column) + operator + b.Wrap(nested))
	return nil
}

func (b *mySQLBuilder) BuildNil(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.Nil)
	if err := b.builder.BuildStatement(stmt, x.Field); err != nil {
		return err
	}
	if x.IsNot {
		stmt.WriteString(" IS NULL")
	} else {
		stmt.WriteString(" IS NOT NULL")
	}
	return nil
}

func (b *mySQLBuilder) BuildRaw(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.Raw)
	stmt.WriteString(x.Value)
	return nil
}

func (b *mySQLBuilder) BuildAs(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.As)
	if err := b.getValue(stmt, x.Field); err != nil {
		return err
	}
	stmt.WriteString(" AS ")
	stmt.WriteString(b.Quote(x.Name))
	return nil
}

func (b *mySQLBuilder) BuildAggregate(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.Aggregate)
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

func (b *mySQLBuilder) BuildOperator(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.Operator)
	stmt.WriteRune(' ')
	stmt.WriteString(operatorMap[x])
	stmt.WriteRune(' ')
	return nil
}

func (b *mySQLBuilder) BuildClause(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.C)
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
	return nil
}

func (b *mySQLBuilder) BuildSort(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.Sort)
	stmt.WriteString(b.Quote(x.Field))
	if x.Order == primitive.Descending {
		stmt.WriteRune(' ')
		stmt.WriteString("DESC")
	}
	return nil
}

func (b *mySQLBuilder) BuildKeyValue(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x := it.(primitive.KV)
	stmt.WriteString(b.Quote(string(x.Field)))
	stmt.WriteString(" = ")
	return b.getValue(stmt, x.Value)
}

func (b *mySQLBuilder) BuildJSONContains(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(primitive.JC)
	stmt.WriteString("JSON_CONTAINS")
	stmt.WriteByte('(')
	if err := b.builder.BuildStatement(stmt, x.Target); err != nil {
		return err
	}
	stmt.WriteByte(',')
	if err := b.builder.BuildStatement(stmt, x.Candidate); err != nil {
		return err
	}
	if x.Path != nil {
		stmt.WriteByte(',')
		stmt.WriteString("'" + *x.Path + "'")
	}
	stmt.WriteByte(')')
	return nil
}

func (b *mySQLBuilder) BuildMath(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x := it.(primitive.Math)
	stmt.WriteString(b.Quote(string(x.Field)) + " ")
	if x.Mode == primitive.Add {
		stmt.WriteRune('+')
	} else {
		stmt.WriteRune('-')
	}
	stmt.WriteString(" " + strconv.Itoa(x.Value))
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
	if builder, ok := b.builder.LookupBuilder(t); ok {
		if err := builder(stmt, it); err != nil {
			return err
		}
		return nil
	}

	stmt.WriteByte('?')
	encoder, err := b.registry.LookupEncoder(v)
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
	x := it.(primitive.Group)
	for len(x.Values) > 0 {
		if err := b.getValue(stmt, x.Values[0]); err != nil {
			return err
		}
		x.Values = x.Values[1:]
	}
	return
}

func (b *mySQLBuilder) BuildRange(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x := it.(primitive.R)
	v := reflext.ValueOf(x.From)
	encoder, err := b.registry.LookupEncoder(v)
	if err != nil {
		return err
	}
	arg, err := encoder(nil, v)
	if err != nil {
		return err
	}
	stmt.AppendArg(arg)

	v = reflext.ValueOf(x.To)
	encoder, err = b.registry.LookupEncoder(v)
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

func (b *mySQLBuilder) BuildFindActions(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(*actions.FindActions)
	x.Table = strings.TrimSpace(x.Table)
	if x.Table == "" {
		return errors.New("mysql: empty table name")
	}
	stmt.WriteString("SELECT ")
	if x.DistinctOn {
		stmt.WriteString("DISTINCT ")
	}
	if err := b.appendSelect(stmt, x.Projections); err != nil {
		return err
	}
	stmt.WriteString(" FROM " + b.TableName(x.Database, x.Table))
	if err := b.appendWhere(stmt, x.Conditions.Values); err != nil {
		return err
	}
	if err := b.appendGroupBy(stmt, x.GroupBys); err != nil {
		return err
	}
	if err := b.appendOrderBy(stmt, x.Sorts); err != nil {
		return err
	}
	b.appendLimitNOffset(stmt, x.Count, x.Skip)

	return nil
}

func (b *mySQLBuilder) BuildUpdateActions(stmt *sqlstmt.Statement, it interface{}) error {
	x, ok := it.(*actions.UpdateActions)
	if !ok {
		return errors.New("data type not match")
	}
	stmt.WriteString("UPDATE " + b.TableName(x.Database, x.Table) + ` `)
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

func (b *mySQLBuilder) BuildDeleteActions(stmt *sqlstmt.Statement, it interface{}) error {
	x := it.(*actions.DeleteActions)
	stmt.WriteString("DELETE FROM " + b.TableName(x.Database, x.Table))
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
	stmt.WriteString("*")
	return nil
}

func (b *mySQLBuilder) appendWhere(stmt *sqlstmt.Statement, conds []interface{}) error {
	length := len(conds)
	if length > 0 {
		stmt.WriteString(" WHERE ")
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
		stmt.WriteString(" GROUP BY ")
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

func (b *mySQLBuilder) appendOrderBy(stmt *sqlstmt.Statement, sorts []interface{}) error {
	length := len(sorts)
	if length < 1 {
		return nil
	}
	stmt.WriteString(" ORDER BY ")
	for i := 0; i < length; i++ {
		if i > 0 {
			stmt.WriteRune(',')
		}
		if err := b.builder.BuildStatement(stmt, sorts[i]); err != nil {
			return err
		}
	}
	return nil
}

func (b *mySQLBuilder) appendLimitNOffset(stmt *sqlstmt.Statement, limit, offset uint) {
	if limit > 0 {
		stmt.WriteString(" LIMIT " + strconv.FormatUint(uint64(limit), 10))
	}
	if offset > 0 {
		stmt.WriteString(" OFFSET " + strconv.FormatUint(uint64(offset), 10))
	}
}

func (b *mySQLBuilder) appendSet(stmt *sqlstmt.Statement, values []primitive.KV) error {
	length := len(values)
	if length > 0 {
		stmt.WriteString("SET ")
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

func escapeWildCard(n string) string {
	length := len(n) - 1
	if length < 1 {
		return n
	}
	blr := new(strings.Builder)
	for i := 0; i < length; i++ {
		switch n[i] {
		case '%':
			blr.WriteString(`\%`)
		case '_':
			blr.WriteString(`\_`)
		case '\\':
			blr.WriteString(`\\`)
		default:
			blr.WriteByte(n[i])
		}
	}
	blr.WriteByte(n[length])
	return blr.String()
}

func unmatchedDataType(callback string) error {
	return errors.New("mysql: invalid data type")
}
