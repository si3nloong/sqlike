package sqlike

import (
	"database/sql"
	"io"
	"reflect"

	"errors"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/sqlike/actions"
)

// ErrNoRows :
var ErrNoRows = sql.ErrNoRows

// EOF :
var EOF = io.EOF

// Result :
type Result struct {
	close    bool
	rows     *sql.Rows
	registry *codec.Registry
	actions  actions.FindActions
	columns  []string
	err      error
}

// Columns :
func (r *Result) Columns() []string {
	return r.columns
}

// ColumnTypes :
func (r *Result) ColumnTypes() ([]*sql.ColumnType, error) {
	return r.rows.ColumnTypes()
}

func (r *Result) nextValues() ([]interface{}, error) {
	if !r.Next() {
		return nil, EOF
	}
	return r.values()
}

func (r *Result) values() ([]interface{}, error) {
	length := len(r.columns)
	values := make([]interface{}, length)
	for j := 0; j < length; j++ {
		values[j] = &values[j]
	}
	if err := r.rows.Scan(values...); err != nil {
		return nil, err
	}
	return values, nil
}

// Scan :
func (r *Result) Scan(dests ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	if len(dests) == 0 {
		return errors.New("sqlike: empty destination to scan")
	}
	values, err := r.values()
	if err != nil {
		return err
	}
	max := len(dests)
	for i, v := range values {
		if i >= max {
			break
		}
		fv := reflext.ValueOf(dests[i])
		if fv.Kind() != reflect.Ptr {
			return ErrUnaddressableEntity
		}
		t := fv.Elem().Type()
		vv := reflext.Zero(t)
		decoder, err := r.registry.LookupDecoder(t)
		if err != nil {
			return err
		}
		if err := decoder(v, vv); err != nil {
			return err
		}
		fv.Elem().Set(vv)
	}
	return nil
}

// Decode will decode the current document into val.
func (r *Result) Decode(dst interface{}) error {
	if r.close {
		defer r.Close()
	}
	if r.err != nil {
		return r.err
	}

	v := reflext.ValueOf(dst)
	if !v.IsValid() {
		return ErrInvalidInput
	}

	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return ErrUnaddressableEntity
	}

	v = reflext.Indirect(v)
	t = reflext.Deref(t)
	if !reflext.IsKind(t, reflect.Struct) {
		return errors.New("sqlike: it must be a struct to decode")
	}

	mapper := reflext.DefaultMapper
	idxs := mapper.TraversalsByName(t, r.columns)
	values, err := r.values()
	if err != nil {
		return err
	}
	vv := reflext.Zero(t)
	for j, idx := range idxs {
		if idx == nil {
			continue
		}
		fv := mapper.FieldByIndexes(vv, idx)
		decoder, err := r.registry.LookupDecoder(fv.Type())
		if err != nil {
			return err
		}
		if err := decoder(values[j], fv); err != nil {
			return err
		}
	}
	reflext.Indirect(v).Set(reflext.Indirect(vv))
	return nil
}

// ScanSlice :
func (r *Result) ScanSlice(results interface{}) error {
	defer r.Close()
	if r.err != nil {
		return r.err
	}

	v := reflext.ValueOf(results)
	if !v.IsValid() {
		return ErrInvalidInput
	}

	if !reflext.IsKind(v.Type(), reflect.Ptr) {
		return ErrUnaddressableEntity
	}

	v = reflext.Indirect(v)
	t := v.Type()
	if !reflext.IsKind(t, reflect.Slice) {
		return errors.New("sqlike: it must be a slice of entity")
	}

	slice := reflect.MakeSlice(t, 0, 0)
	t = t.Elem()
	// decoders := make([]coder.ValueDecoder, length, length)

	for i := 0; r.rows.Next(); i++ {
		values, err := r.values()
		if err != nil {
			return err
		}
		slice = reflect.Append(slice, reflext.Zero(t))
		fv := slice.Index(i)
		decoder, err := r.registry.LookupDecoder(fv.Type())
		if err != nil {
			return err
		}
		if err := decoder(values[0], fv); err != nil {
			return err
		}
	}
	v.Set(slice)
	return nil
}

// All :
func (r *Result) All(results interface{}) error {
	defer r.Close()
	if r.err != nil {
		return r.err
	}

	v := reflext.ValueOf(results)
	if !v.IsValid() {
		return ErrInvalidInput
	}

	if !reflext.IsKind(v.Type(), reflect.Ptr) {
		return ErrUnaddressableEntity
	}

	v = reflext.Indirect(v)
	t := v.Type()
	if !reflext.IsKind(t, reflect.Slice) {
		return errors.New("sqlike: it must be a slice of entity")
	}

	length := len(r.columns)
	slice := reflect.MakeSlice(t, 0, 0)
	t = t.Elem()
	mapper := reflext.DefaultMapper
	idxs := mapper.TraversalsByName(t, r.columns)
	decoders := make([]codec.ValueDecoder, length)
	for i := 0; r.rows.Next(); i++ {
		values, err := r.values()
		if err != nil {
			return err
		}
		vv := reflext.Zero(t)
		for j, idx := range idxs {
			if idx == nil {
				continue
			}
			fv := mapper.FieldByIndexes(vv, idx)
			if i < 1 {
				decoder, err := r.registry.LookupDecoder(fv.Type())
				if err != nil {
					return err
				}
				decoders[j] = decoder
			}
			if err := decoders[j](values[j], fv); err != nil {
				return err
			}
		}
		slice = reflect.Append(slice, vv)
	}
	v.Set(slice)
	return nil
}

// Error :
func (r *Result) Error() error {
	if r.rows != nil {
		defer r.rows.Close()
	}
	return r.err
}

// Next :
func (r *Result) Next() bool {
	return r.rows.Next()
}

// Close :
func (r *Result) Close() error {
	if r.rows != nil {
		return r.rows.Close()
	}
	return nil
}
