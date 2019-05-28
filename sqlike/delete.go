package sqlike

import (
	"reflect"
	"bitbucket.org/SianLoong/sqlike/core"
	"bitbucket.org/SianLoong/sqlike/reflext"
	"bitbucket.org/SianLoong/sqlike/sqlike/actions"
	sqldriver "bitbucket.org/SianLoong/sqlike/sqlike/sql/driver"
	"bitbucket.org/SianLoong/sqlike/sqlike/sql/expr"
	"golang.org/x/xerrors"
)

// DestroyOne :
func (tb *Table) DestroyOne(delete interface{}) error {
	v := reflect.ValueOf(delete)
	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return ErrUnaddressableEntity
	}

	if v.IsNil() {
		return xerrors.New("entity is nil")
	}

	mapper := core.DefaultMapper
	cdc := mapper.CodecByType(t)
	f, exists := cdc.Names["$Key"]
	if !exists {
		return xerrors.New(`missing $Key field`)
	}

	x := new(actions.DeleteActions)
	x.Table = tb.name
	fv := mapper.FieldByIndexesReadOnly(v, f.Index)
	x.Where(expr.Equal(f.Path, fv.Interface()))
	x.Limit(1)
	stmt, err := tb.dialect.Delete(x)
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

// DeleteMany :
func (tb *Table) DeleteMany(act actions.DeleteStatement) (int64, error) {
	x := new(actions.DeleteActions)
	if act != nil {
		*x = *(act.(*actions.DeleteActions))
	}
	if x.Table == "" {
		x.Table = tb.name
	}

	if len(x.Conditions) < 1 {
		return 0, xerrors.New("sqlike: no condition is not allow for delete")
	}

	stmt, err := tb.dialect.Delete(x)
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
