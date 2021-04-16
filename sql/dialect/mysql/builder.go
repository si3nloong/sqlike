package mysql

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/spatial"
	"github.com/si3nloong/sqlike/sql"
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
	registry codec.Codecer
	builder  *sqlstmt.StatementBuilder
	sqlutil.MySQLUtil
}

func (b mySQLBuilder) SetRegistryAndBuilders(rg codec.Codecer, blr *sqlstmt.StatementBuilder) {
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
	blr.SetBuilder(reflect.TypeOf(primitive.Encoding{}), b.BuildEncoding)
	blr.SetBuilder(reflect.TypeOf(primitive.Aggregate{}), b.BuildAggregate)
	blr.SetBuilder(reflect.TypeOf(primitive.Column{}), b.BuildColumn)
	blr.SetBuilder(reflect.TypeOf(primitive.JSONColumn{}), b.BuildJSONColumn)
	blr.SetBuilder(reflect.TypeOf(primitive.C{}), b.BuildClause)
	blr.SetBuilder(reflect.TypeOf(primitive.L{}), b.BuildLike)
	blr.SetBuilder(reflect.TypeOf(primitive.TypeSafe{}), b.BuildTypeSafe)
	blr.SetBuilder(reflect.TypeOf(primitive.Operator(0)), b.BuildOperator)
	blr.SetBuilder(reflect.TypeOf(primitive.Group{}), b.BuildGroup)
	blr.SetBuilder(reflect.TypeOf(primitive.R{}), b.BuildRange)
	blr.SetBuilder(reflect.TypeOf(primitive.Sort{}), b.BuildSort)
	blr.SetBuilder(reflect.TypeOf(primitive.KV{}), b.BuildKeyValue)
	blr.SetBuilder(reflect.TypeOf(primitive.Math{}), b.BuildMath)
	blr.SetBuilder(reflect.TypeOf(&primitive.Case{}), b.BuildCase)
	blr.SetBuilder(reflect.TypeOf(spatial.Func{}), b.BuildSpatialFunc)
	blr.SetBuilder(reflect.TypeOf(&sql.SelectStmt{}), b.BuildSelectStmt)
	blr.SetBuilder(reflect.TypeOf(&sql.UpdateStmt{}), b.BuildUpdateStmt)
	// blr.SetBuilder(reflect.TypeOf(&sql.DeleteStmt{}), b.BuildDeleteStmt)
	blr.SetBuilder(reflect.TypeOf(&actions.FindActions{}), b.BuildFindActions)
	blr.SetBuilder(reflect.TypeOf(&actions.UpdateActions{}), b.BuildUpdateActions)
	blr.SetBuilder(reflect.TypeOf(&actions.DeleteActions{}), b.BuildDeleteActions)
	blr.SetBuilder(reflect.String, b.BuildString)
	b.registry = rg
	b.builder = blr
}

// BuildCastAs :
func (b *mySQLBuilder) BuildCastAs(stmt sqlstmt.Stmt, it interface{}) error {
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

// BuildFunction :
func (b *mySQLBuilder) BuildFunction(stmt sqlstmt.Stmt, it interface{}) error {
	x := it.(primitive.Func)
	stmt.WriteString(x.Name)
	stmt.WriteByte('(')
	for i, args := range x.Args {
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

// BuildJSONFunction :
func (b *mySQLBuilder) BuildJSONFunction(stmt sqlstmt.Stmt, it interface{}) error {
	x := it.(primitive.JSONFunc)
	stmt.WriteString(x.Type.String())
	stmt.WriteByte('(')
	for i, args := range x.Args {
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

// BuildString :
func (b *mySQLBuilder) BuildString(stmt sqlstmt.Stmt, it interface{}) error {
	v := reflect.ValueOf(it)
	stmt.WriteString(b.Quote(v.String()))
	return nil
}

// BuildLike :
func (b *mySQLBuilder) BuildLike(stmt sqlstmt.Stmt, it interface{}) error {
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
		stmt.AppendArgs(nil)
		return nil
	}

	t := v.Type()
	if builder, ok := b.builder.LookupBuilder(t); ok {
		if err := builder(stmt, x.Value); err != nil {
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
	stmt.AppendArgs(vv)
	return nil
}

// BuildField :
func (b *mySQLBuilder) BuildField(stmt sqlstmt.Stmt, it interface{}) error {
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

// BuildValue :
func (b *mySQLBuilder) BuildValue(stmt sqlstmt.Stmt, it interface{}) (err error) {
	x := it.(primitive.Value)
	v := reflext.ValueOf(x.Raw)
	if !v.IsValid() {
		stmt.WriteByte('?')
		stmt.AppendArgs(nil)
		return
	}

	encoder, err := b.registry.LookupEncoder(v)
	if err != nil {
		return err
	}
	vv, err := encoder(nil, v)
	if err != nil {
		return err
	}
	convertSpatial(stmt, vv)
	// stmt.WriteRune('?')
	// stmt.AppendArgs(vv)
	return nil
}

// BuildColumn :
func (b *mySQLBuilder) BuildColumn(stmt sqlstmt.Stmt, it interface{}) error {
	x := it.(primitive.Column)
	if x.Table != "" {
		stmt.WriteString(b.Quote(x.Table))
		stmt.WriteByte('.')
	}
	stmt.WriteString(b.Quote(x.Name))
	return nil
}

// BuildJSONColumn :
func (b *mySQLBuilder) BuildJSONColumn(stmt sqlstmt.Stmt, it interface{}) error {
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

// BuildNil :
func (b *mySQLBuilder) BuildNil(stmt sqlstmt.Stmt, it interface{}) error {
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

// BuildRaw :
func (b *mySQLBuilder) BuildRaw(stmt sqlstmt.Stmt, it interface{}) error {
	x, ok := it.(primitive.Raw)
	if ok {
		stmt.WriteString(x.Value)
	}
	return nil
}

// BuildAs :
func (b *mySQLBuilder) BuildAs(stmt sqlstmt.Stmt, it interface{}) error {
	x := it.(primitive.As)
	if err := b.getValue(stmt, x.Field); err != nil {
		return err
	}
	stmt.WriteString(" AS ")
	stmt.WriteString(b.Quote(x.Name))
	return nil
}

// BuildAggregate :
func (b *mySQLBuilder) BuildAggregate(stmt sqlstmt.Stmt, it interface{}) error {
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

// BuildOperator :
func (b *mySQLBuilder) BuildOperator(stmt sqlstmt.Stmt, it interface{}) error {
	x := it.(primitive.Operator)
	stmt.WriteByte(' ')
	stmt.WriteString(operatorMap[x])
	stmt.WriteByte(' ')
	return nil
}

// BuildClause :
func (b *mySQLBuilder) BuildClause(stmt sqlstmt.Stmt, it interface{}) error {
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

// BuildSort :
func (b *mySQLBuilder) BuildSort(stmt sqlstmt.Stmt, it interface{}) error {
	x := it.(primitive.Sort)
	if err := b.builder.BuildStatement(stmt, x.Field); err != nil {
		return err
	}
	if x.Order == primitive.Descending {
		stmt.WriteByte(' ')
		stmt.WriteString("DESC")
	}
	return nil
}

// BuildKeyValue :
func (b *mySQLBuilder) BuildKeyValue(stmt sqlstmt.Stmt, it interface{}) (err error) {
	x := it.(primitive.KV)
	stmt.WriteString(b.Quote(string(x.Field)))
	stmt.WriteString(" = ")
	return b.getValue(stmt, x.Value)
}

// BuildMath :
func (b *mySQLBuilder) BuildMath(stmt sqlstmt.Stmt, it interface{}) (err error) {
	x := it.(primitive.Math)
	stmt.WriteString(b.Quote(string(x.Field)) + " ")
	if x.Mode == primitive.Add {
		stmt.WriteByte('+')
	} else {
		stmt.WriteByte('-')
	}
	stmt.WriteString(" " + strconv.Itoa(x.Value))
	return
}

// BuildCase :
func (b *mySQLBuilder) BuildCase(stmt sqlstmt.Stmt, it interface{}) error {
	x := it.(*primitive.Case)
	stmt.WriteByte('(')
	stmt.WriteString("CASE")
	for _, w := range x.WhenClauses {
		stmt.WriteString(" WHEN ")
		if err := b.builder.BuildStatement(stmt, w[0]); err != nil {
			return err
		}
		stmt.WriteString(" THEN ")
		if err := b.getValue(stmt, w[1]); err != nil {
			return err
		}
	}
	stmt.WriteString(" ELSE ")
	if x.ElseClause != nil {
		if err := b.getValue(stmt, x.ElseClause); err != nil {
			return err
		}
	}
	stmt.WriteString(" END")
	stmt.WriteByte(')')
	return nil
}

// BuildSpatialFunc :
func (b *mySQLBuilder) BuildSpatialFunc(stmt sqlstmt.Stmt, it interface{}) (err error) {
	x := it.(spatial.Func)
	stmt.WriteString(x.Type.String())
	stmt.WriteByte('(')
	for i, arg := range x.Args {
		if i > 0 {
			stmt.WriteByte(',')
		}
		if err := b.builder.BuildStatement(stmt, arg); err != nil {
			return err
		}
	}
	stmt.WriteByte(')')
	return
}

// BuildGroup :
func (b *mySQLBuilder) BuildGroup(stmt sqlstmt.Stmt, it interface{}) (err error) {
	x := it.(primitive.Group)
	for len(x.Values) > 0 {
		if err := b.getValue(stmt, x.Values[0]); err != nil {
			return err
		}
		x.Values = x.Values[1:]
	}
	return
}

// BuildRange :
func (b *mySQLBuilder) BuildRange(stmt sqlstmt.Stmt, it interface{}) (err error) {
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
	stmt.AppendArgs(arg)

	v = reflext.ValueOf(x.To)
	encoder, err = b.registry.LookupEncoder(v)
	if err != nil {
		return err
	}
	arg, err = encoder(nil, v)
	if err != nil {
		return err
	}
	stmt.AppendArgs(arg)
	stmt.WriteByte('?')
	stmt.WriteString(" AND ")
	stmt.WriteByte('?')
	return
}

// BuildEncoding :
func (b *mySQLBuilder) BuildEncoding(stmt sqlstmt.Stmt, it interface{}) (err error) {
	x := it.(primitive.Encoding)
	if x.Charset != nil {
		if (*x.Charset)[0] != '_' {
			stmt.WriteString("_")
		}
		stmt.WriteString(*x.Charset + " ")
	}
	err = b.builder.BuildStatement(stmt, x.Column)
	if err != nil {
		return
	}
	stmt.WriteString(" COLLATE " + x.Collate)
	return
}

// BuildTypeSafe :
func (b *mySQLBuilder) BuildTypeSafe(stmt sqlstmt.Stmt, it interface{}) (err error) {
	ts := it.(primitive.TypeSafe)
	switch ts.Type {
	case primitive.Varchar:
		stmt.WriteString(strconv.Quote(ts.Value.(string)))
	}
	return
}

// BuildSelectStmt :
func (b *mySQLBuilder) BuildSelectStmt(stmt sqlstmt.Stmt, it interface{}) error {
	x := it.(*sql.SelectStmt)
	stmt.WriteString("SELECT ")
	if x.DistinctOn {
		stmt.WriteString("DISTINCT ")
	}
	if err := b.appendSelect(stmt, x.Projections); err != nil {
		return err
	}
	stmt.WriteString(" FROM ")
	if err := b.appendTable(stmt, x.Tables); err != nil {
		return err
	}
	if err := b.appendWhere(stmt, x.Conditions.Values); err != nil {
		return err
	}
	if err := b.appendGroupBy(stmt, x.Groups); err != nil {
		return err
	}
	if err := b.appendOrderBy(stmt, x.Sorts); err != nil {
		return err
	}
	b.appendLimitNOffset(stmt, x.Max, x.Skip)
	return nil
}

// BuildUpdateStmt :
func (b *mySQLBuilder) BuildUpdateStmt(stmt sqlstmt.Stmt, it interface{}) error {
	x := it.(*sql.UpdateStmt)
	stmt.WriteString("UPDATE " + b.TableName(x.Database, x.Table) + ` `)
	if err := b.appendSet(stmt, x.Values); err != nil {
		return err
	}
	if err := b.appendWhere(stmt, x.Conditions.Values); err != nil {
		return err
	}
	if err := b.appendOrderBy(stmt, x.Sorts); err != nil {
		return err
	}
	b.appendLimitNOffset(stmt, x.Max, 0)
	return nil
}

// BuildFindActions :
func (b *mySQLBuilder) BuildFindActions(stmt sqlstmt.Stmt, it interface{}) error {
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

// BuildUpdateActions :
func (b *mySQLBuilder) BuildUpdateActions(stmt sqlstmt.Stmt, it interface{}) error {
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

// BuildDeleteActions :
func (b *mySQLBuilder) BuildDeleteActions(stmt sqlstmt.Stmt, it interface{}) error {
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

func (b *mySQLBuilder) getValue(stmt sqlstmt.Stmt, it interface{}) (err error) {
	v := reflext.ValueOf(it)
	if !v.IsValid() {
		stmt.WriteByte('?')
		stmt.AppendArgs(nil)
		return
	}

	t := v.Type()
	if builder, ok := b.builder.LookupBuilder(t); ok {
		if err := builder(stmt, it); err != nil {
			return err
		}
		return nil
	}

	encoder, err := b.registry.LookupEncoder(v)
	if err != nil {
		return err
	}
	vv, err := encoder(nil, v)
	if err != nil {
		return err
	}
	// stmt.WriteByte('?')
	// stmt.AppendArgs(vv)
	convertSpatial(stmt, vv)
	return
}

func (b *mySQLBuilder) appendSelect(stmt sqlstmt.Stmt, pjs []interface{}) error {
	if len(pjs) > 0 {
		length := len(pjs)
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteByte(',')
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

func (b *mySQLBuilder) appendTable(stmt sqlstmt.Stmt, fields []interface{}) error {
	length := len(fields)
	if length > 0 {
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteByte(' ')
			}
			if err := b.builder.BuildStatement(stmt, fields[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *mySQLBuilder) appendWhere(stmt sqlstmt.Stmt, conds []interface{}) error {
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

func (b *mySQLBuilder) appendGroupBy(stmt sqlstmt.Stmt, fields []interface{}) error {
	length := len(fields)
	if length > 0 {
		stmt.WriteString(" GROUP BY ")
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteByte(',')
			}
			if err := b.builder.BuildStatement(stmt, fields[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *mySQLBuilder) appendOrderBy(stmt sqlstmt.Stmt, sorts []interface{}) error {
	length := len(sorts)
	if length < 1 {
		return nil
	}
	stmt.WriteString(" ORDER BY ")
	for i := 0; i < length; i++ {
		if i > 0 {
			stmt.WriteByte(',')
		}
		if err := b.builder.BuildStatement(stmt, sorts[i]); err != nil {
			return err
		}
	}
	return nil
}

func (b *mySQLBuilder) appendLimitNOffset(stmt sqlstmt.Stmt, limit, offset uint) {
	if limit > 0 {
		stmt.WriteString(" LIMIT " + strconv.FormatUint(uint64(limit), 10))
	}
	if offset > 0 {
		stmt.WriteString(" OFFSET " + strconv.FormatUint(uint64(offset), 10))
	}
}

func (b *mySQLBuilder) appendSet(stmt sqlstmt.Stmt, values []primitive.KV) error {
	length := len(values)
	if length > 0 {
		stmt.WriteString("SET ")
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteByte(',')
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
