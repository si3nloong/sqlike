package sqlike

import (
	"context"
	"database/sql"
	"io"
	"reflect"

	"errors"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/sql/dialect"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// ErrNoRows : is an alias for no record found
var ErrNoRows = sql.ErrNoRows

// EOF : is an alias for end of file
var EOF = io.EOF

// Resulter :
type Resulter interface {
	Scan(dests ...interface{}) error
	Columns() []string
	Next() bool
	NextResultSet() bool
	Close() error
}

// Result :
type Result struct {
	ctx         context.Context
	close       bool
	rows        *sql.Rows
	cache       reflext.StructMapper
	dialect     dialect.Dialect
	columns     []string
	columnTypes []*sql.ColumnType
	err         error
}

var _ Resulter = (*Result)(nil)

// Columns :
func (r *Result) Columns() []string {
	return r.columns
}

// ColumnTypes :
func (r *Result) ColumnTypes() ([]*sql.ColumnType, error) {
	return r.columnTypes, nil
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

// Scan : will behave as similar as sql.Scan.
func (r *Result) Scan(dests ...interface{}) error {
	if r.close {
		defer r.Close()
	}
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
		fv = reflext.IndirectInit(fv)
		decoder, err := r.dialect.LookupDecoder(fv.Type())
		if err != nil {
			return err
		}
		if err := decoder(r.ctx, v, fv); err != nil {
			return err
		}
	}
	return nil
}

// Decode will decode the current document into val, this will only accepting pointer of struct as an input.
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

	idxs := r.cache.TraversalsByName(t, r.columns)
	values, err := r.values()
	if err != nil {
		return err
	}
	vv := reflext.Zero(t)
	for j, idx := range idxs {
		if idx == nil {
			continue
		}
		fv := r.cache.FieldByIndexes(vv, idx)
		decoder, err := r.dialect.LookupDecoder(fv.Type())
		if err != nil {
			return err
		}
		if err := decoder(r.ctx, values[j], fv); err != nil {
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

	for i := 0; r.rows.Next(); i++ {
		values, err := r.values()
		if err != nil {
			return err
		}
		slice = reflect.Append(slice, reflext.Zero(t))
		fv := slice.Index(i)
		decoder, err := r.dialect.LookupDecoder(fv.Type())
		if err != nil {
			return err
		}
		if err := decoder(r.ctx, values[0], fv); err != nil {
			return err
		}
	}
	v.Set(slice)
	return r.rows.Close()
}

// All : this will map all the records from sql to a slice of struct.
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
	idxs := r.cache.TraversalsByName(t, r.columns)
	decoders := make([]db.ValueDecoder, length)
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
			fv := r.cache.FieldByIndexes(vv, idx)
			if i < 1 {
				decoder, err := r.dialect.LookupDecoder(fv.Type())
				if err != nil {
					return err
				}
				decoders[j] = decoder
			}
			if err := decoders[j](r.ctx, values[j], fv); err != nil {
				return err
			}
		}
		slice = reflect.Append(slice, vv)
	}
	v.Set(slice)
	return r.rows.Close()
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

// NextResultSet :
func (r *Result) NextResultSet() bool {
	return r.rows.NextResultSet()
}

// Close :
func (r *Result) Close() error {
	if r.rows != nil {
		return r.rows.Close()
	}
	return nil
}
