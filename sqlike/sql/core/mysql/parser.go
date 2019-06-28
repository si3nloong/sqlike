package mysql

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/si3nloong/sqlike/sqlike/sql/codec"
	sqlstmt "github.com/si3nloong/sqlike/sqlike/sql/stmt"
	sqlutil "github.com/si3nloong/sqlike/sqlike/sql/util"
	"golang.org/x/xerrors"
)

var operatorMap = map[primitive.Operator]string{
	primitive.Equal:        "=",
	primitive.NotEqual:     "<>",
	primitive.Like:         "LIKE",
	primitive.NotLike:      "NOT LIKE",
	primitive.In:           "IN",
	primitive.NotIn:        "NOT IN",
	primitive.Between:      "BETWEEN",
	primitive.NotBetween:   "NOT BETWEEN",
	primitive.IsNull:       "IS NULL",
	primitive.NotNull:      "IS NOT NULL",
	primitive.GreaterThan:  ">",
	primitive.GreaterEqual: ">=",
	primitive.LowerThan:    "<",
	primitive.LowerEqual:   "<=",
	primitive.Or:           "OR",
	primitive.And:          "AND",
}

type mySQLParser struct {
	registry *codec.Registry
	parser   *sqlstmt.StatementParser
	sqlutil.MySQLUtil
}

func (p mySQLParser) SetParsers(parser *sqlstmt.StatementParser) {
	parser.SetParser(reflect.TypeOf(primitive.L("")), p.ParseString)
	parser.SetParser(reflect.TypeOf(primitive.Raw("")), p.ParseRaw)
	parser.SetParser(reflect.TypeOf(primitive.C{}), p.ParseClause)
	parser.SetParser(reflect.TypeOf(primitive.Col("")), p.ParseString)
	parser.SetParser(reflect.TypeOf(primitive.Operator(0)), p.ParseOperator)
	parser.SetParser(reflect.TypeOf(primitive.G{}), p.ParseGroup)
	parser.SetParser(reflect.TypeOf(primitive.GV{}), p.ParseGroupValue)
	parser.SetParser(reflect.TypeOf(primitive.R{}), p.ParseRange)
	parser.SetParser(reflect.TypeOf(primitive.Sort{}), p.ParseSort)
	parser.SetParser(reflect.TypeOf(primitive.KV{}), p.ParseKeyValue)
	parser.SetParser(reflect.TypeOf(primitive.Math{}), p.ParseMath)
	parser.SetParser(reflect.TypeOf(&actions.FindActions{}), p.ParseFindActions)
	parser.SetParser(reflect.TypeOf(&actions.UpdateActions{}), p.ParseUpdateActions)
	parser.SetParser(reflect.TypeOf(&actions.DeleteActions{}), p.ParseDeleteActions)
	parser.SetParser(reflect.String, p.ParseString)
	p.parser = parser
}

func (p mySQLParser) ParseString(stmt *sqlstmt.Statement, it interface{}) error {
	v := reflect.ValueOf(it)
	stmt.WriteString(p.Quote(v.String()))
	return nil
}

func (p mySQLParser) ParseRaw(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.Raw)
	if !isOk {
		return xerrors.New("raw")
	}
	stmt.WriteString(string(x))
	return nil
}

func (p mySQLParser) ParseOperator(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.Operator)
	if !isOk {
		return xerrors.New("operator")
	}
	stmt.WriteRune(' ')
	stmt.WriteString(operatorMap[x])
	stmt.WriteRune(' ')
	return nil
}

func (p *mySQLParser) ParseClause(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.C)
	if !isOk {
		return xerrors.New("clause")
	}

	if err := p.parser.BuildStatement(stmt, x.Field); err != nil {
		return err
	}

	stmt.WriteString(" " + operatorMap[x.Operator] + " ")
	switch x.Operator {
	case primitive.IsNull, primitive.NotNull:
		return nil
	}

	if x.Value == nil {
		stmt.AppendArg(nil)
		stmt.WriteRune('?')
		return nil
	}

	if parser, isOk := p.parser.LookupParser(x.Value); isOk {
		stmt.WriteRune('(')
		if err := parser(stmt, x.Value); err != nil {
			return err
		}
		stmt.WriteRune(')')
		return nil
	}

	v := reflext.ValueOf(x.Value)
	encoder, err := codec.DefaultRegistry.LookupEncoder(v.Type())
	if err != nil {
		return err
	}
	arg, err := encoder(nil, v)
	if err != nil {
		return err
	}
	stmt.AppendArg(arg)
	stmt.WriteRune('?')
	return nil
}

func (p *mySQLParser) ParseSort(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(primitive.Sort)
	if !isOk {
		return xerrors.New("sort")
	}
	stmt.WriteString(p.Quote(x.Field))
	if x.Order == primitive.Descending {
		stmt.WriteRune(' ')
		stmt.WriteString(`DESC`)
	}
	return nil
}

func (p *mySQLParser) ParseKeyValue(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x, isOk := it.(primitive.KV)
	if !isOk {
		return xerrors.New("data type not match")
	}
	stmt.WriteString(p.Quote(string(x.Field)))
	stmt.WriteString(` = `)
	var parser sqlstmt.ParseStatementFunc
	parser, isOk = p.parser.LookupParser(reflext.ValueOf(it))
	if isOk {
		err = parser(stmt, it)
		return
	}
	// log.Println(parser)
	stmt.WriteString(`?`)
	// log.Println("Value should normalize before add into arguments")
	// log.Println("Value :", x.Value)
	stmt.AppendArg(x.Value)
	return
}

func (p *mySQLParser) ParseMath(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x, isOk := it.(primitive.Math)
	if !isOk {
		return xerrors.New("data type not match")
	}
	stmt.WriteString(p.Quote(string(x.Field)) + ` `)
	if x.Mode == primitive.Add {
		stmt.WriteRune('+')
	} else {
		stmt.WriteRune('-')
	}
	stmt.WriteString(` ` + strconv.Itoa(x.Value))
	return
}

func (p *mySQLParser) ParseGroup(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x, isOk := it.(primitive.G)
	if !isOk {
		return xerrors.New("data type not match")
	}
	for len(x) > 0 {
		err = p.parser.BuildStatement(stmt, x[0])
		if err != nil {
			return
		}
		x = x[1:]
	}
	return
}

func (p *mySQLParser) ParseGroupValue(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x, isOk := it.(primitive.GV)
	if !isOk {
		return xerrors.New("expected data type primitive.GV")
	}
	for len(x) > 0 {
		v := reflext.ValueOf(x[0])
		encoder, err := codec.DefaultRegistry.LookupEncoder(v.Type())
		if err != nil {
			return err
		}
		arg, err := encoder(nil, v)
		if err != nil {
			return err
		}
		stmt.AppendArg(arg)
		stmt.WriteRune('?')
		x = x[1:]
		if len(x) > 0 {
			stmt.WriteRune(',')
		}
	}
	return
}

func (p *mySQLParser) ParseRange(stmt *sqlstmt.Statement, it interface{}) (err error) {
	x, isOk := it.(primitive.R)
	if !isOk {
		return xerrors.New("expected data type primitive.GV")
	}

	v := reflext.ValueOf(x.From)
	encoder, err := codec.DefaultRegistry.LookupEncoder(v.Type())
	if err != nil {
		return err
	}
	arg, err := encoder(nil, v)
	if err != nil {
		return err
	}
	stmt.AppendArg(arg)

	v = reflext.ValueOf(x.To)
	encoder, err = codec.DefaultRegistry.LookupEncoder(v.Type())
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

func (p *mySQLParser) ParseFindActions(stmt *sqlstmt.Statement, it interface{}) error {
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
	if err := p.appendSelect(stmt, x.Projections); err != nil {
		return err
	}
	stmt.WriteString(` FROM ` + p.Quote(x.Table))
	if err := p.appendWhere(stmt, x.Conditions); err != nil {
		return err
	}
	if err := p.appendGroupBy(stmt, x.GroupBys); err != nil {
		return err
	}
	if err := p.appendOrderBy(stmt, x.Sorts); err != nil {
		return err
	}
	p.appendLimitNOffset(stmt, x.Record, x.Skip)

	return nil
}

func (p *mySQLParser) ParseUpdateActions(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(*actions.UpdateActions)
	if !isOk {
		return xerrors.New("data type not match")
	}
	stmt.WriteString(`UPDATE ` + p.Quote(x.Table) + ` `)
	if err := p.appendSet(stmt, x.Values); err != nil {
		return err
	}
	if err := p.appendWhere(stmt, x.Conditions); err != nil {
		return err
	}
	if err := p.appendOrderBy(stmt, x.Sorts); err != nil {
		return err
	}
	p.appendLimitNOffset(stmt, x.Record, 0)
	return nil
}

func (p *mySQLParser) ParseDeleteActions(stmt *sqlstmt.Statement, it interface{}) error {
	x, isOk := it.(*actions.DeleteActions)
	if !isOk {
		return xerrors.New("data type not match")
	}
	stmt.WriteString(`DELETE FROM ` + p.Quote(x.Table))
	if err := p.appendWhere(stmt, x.Conditions); err != nil {
		return err
	}
	if err := p.appendOrderBy(stmt, x.Sorts); err != nil {
		return err
	}
	p.appendLimitNOffset(stmt, x.Record, 0)
	return nil
}

func (p *mySQLParser) appendSelect(stmt *sqlstmt.Statement, pjs []interface{}) error {
	if len(pjs) > 0 {
		length := len(pjs)
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteRune(',')
			}
			if err := p.parser.BuildStatement(stmt, pjs[i]); err != nil {
				return err
			}
		}
		return nil
	}

	stmt.WriteString(`*`)
	return nil
}

func (p *mySQLParser) appendWhere(stmt *sqlstmt.Statement, conds []interface{}) error {
	length := len(conds)
	if length > 0 {
		stmt.WriteString(` WHERE `)
		for i := 0; i < length; i++ {
			if err := p.parser.BuildStatement(stmt, conds[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *mySQLParser) appendGroupBy(stmt *sqlstmt.Statement, fields []interface{}) error {
	length := len(fields)
	if length > 0 {
		stmt.WriteString(` GROUP BY `)
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteRune(',')
			}
			if err := p.parser.BuildStatement(stmt, fields[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *mySQLParser) appendOrderBy(stmt *sqlstmt.Statement, sorts []primitive.Sort) error {
	length := len(sorts)
	if length > 0 {
		stmt.WriteString(` ORDER BY `)
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteRune(',')
			}
			if err := p.parser.BuildStatement(stmt, sorts[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *mySQLParser) appendLimitNOffset(stmt *sqlstmt.Statement, limit, offset uint) {
	if limit > 0 {
		stmt.WriteString(` LIMIT ` + strconv.FormatUint(uint64(limit), 10))
	}
	if offset > 0 {
		stmt.WriteString(` OFFSET ` + strconv.FormatUint(uint64(offset), 10))
	}
}

func (p *mySQLParser) appendSet(stmt *sqlstmt.Statement, values []primitive.C) error {
	length := len(values)
	if length > 0 {
		stmt.WriteString(`SET `)
		for i := 0; i < length; i++ {
			if i > 0 {
				stmt.WriteRune(',')
			}
			if err := p.parser.BuildStatement(stmt, values[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
