package sqlike

import (
	"reflect"
	"github.com/si3nloong/sqlike/core"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	sqldriver "github.com/si3nloong/sqlike/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
	"golang.org/x/xerrors"
)

// ModifyOne :
func (tb *Table) ModifyOne(update interface{}, opts ...*options.ModifyOneOptions) error {
	v := reflect.ValueOf(update)
	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return ErrUnaddressableEntity
	}

	if v.IsNil() {
		return xerrors.New("entity is nil")
	}

	mapper := core.DefaultMapper
	cdc := mapper.CodecByType(t)
	if _, exists := cdc.Names["$Key"]; !exists {
		return xerrors.New(`missing $Key field`)
	}

	_, fields := skipGeneratedColumns(cdc.NameFields)
	x := new(actions.UpdateActions)
	x.Table = tb.name

	for _, sf := range fields {
		fv := mapper.FieldByIndexesReadOnly(v, sf.Index)
		if sf.Path == "$Key" {
			x.Where(expr.Equal("$Key", fv.Interface()))
			continue
		}
		x.Set(sf.Path, fv.Interface())
	}

	x.Limit(1)
	stmt, err := tb.dialect.Update(x)
	if err != nil {
		return err
	}

	result, err := sqldriver.Execute(
		tb.driver,
		stmt,
		tb.logger,
	)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected <= 0 {
		return xerrors.New("unable to modify entity")
	}
	return err
}

// UpdateOne :
func (tb *Table) UpdateOne(update interface{}, opts ...*options.UpdateOneOptions) (int64, error) {
	v := reflect.ValueOf(update)
	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return 0, ErrUnaddressableEntity
	}

	if v.IsNil() {
		return 0, xerrors.New("entity is nil")
	}

	x := new(actions.UpdateActions)
	x.Table = tb.name
	mapper := core.DefaultMapper
	cdc := mapper.CodecByType(t)
	_, fields := skipGeneratedColumns(cdc.NameFields)

	for _, sf := range fields {
		fv := mapper.FieldByIndexesReadOnly(v, sf.Index)
		if sf.Path == "$Key" {
			x.Where(expr.Equal("$Key", fv.Interface()))
			continue
		}
		x.Set(sf.Path, fv.Interface())
	}

	x.Limit(1)
	stmt, err := tb.dialect.Update(x)
	if err != nil {
		return 0, err
	}

	result, err := sqldriver.Execute(
		tb.driver,
		stmt,
		tb.logger,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateMany :
func (tb *Table) UpdateMany(act actions.UpdateStatement) (int64, error) {
	x := new(actions.UpdateActions)
	if act != nil {
		*x = *(act.(*actions.UpdateActions))
	}
	if x.Table == "" {
		x.Table = tb.name
	}

	if len(x.Values) < 1 {
		return 0, xerrors.New("sqlike: no value to update")
	}

	stmt, err := tb.dialect.Update(x)
	if err != nil {
		return 0, err
	}

	result, err := sqldriver.Execute(
		tb.driver,
		stmt,
		tb.logger,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
