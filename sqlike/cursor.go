package sqlike

import (
	"database/sql"
	"reflect"

	"github.com/si3nloong/sqlike/core"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/sql/codec"
	"golang.org/x/xerrors"
)

// Cursor :
type Cursor struct {
	rows     *sql.Rows
	registry *codec.Registry
	columns  []string
	err      error
}

// Columns :
func (c *Cursor) Columns() []string {
	return c.columns
}

// ColumnTypes :
func (c *Cursor) ColumnTypes() ([]*sql.ColumnType, error) {
	return c.rows.ColumnTypes()
}

// Decode will decode the current document into val.
func (c *Cursor) Decode(dst interface{}) error {
	if c.err != nil {
		return c.err
	}

	v := reflext.ValueOf(dst)
	t := v.Type()
	if !reflext.IsKind(t, reflect.Ptr) {
		return ErrUnaddressableEntity
	}

	v = reflext.Indirect(v)
	t = reflext.Deref(t)
	if !reflext.IsKind(t, reflect.Struct) {
		return xerrors.New("it must be struct to decode")
	}

	length := len(c.columns)
	mapper := core.DefaultMapper
	idxs := mapper.TraversalsByName(t, c.columns)
	values := make([]interface{}, length, length)

	for j := 0; j < length; j++ {
		values[j] = &values[j]
	}
	if err := c.rows.Scan(values...); err != nil {
		return err
	}

	vv := reflext.Zero(t)
	for j, idx := range idxs {
		fv := mapper.FieldByIndexes(vv, idx)
		// log.Println(j, fv.Type())
		decoder, err := c.registry.LookupDecoder(fv.Type())
		if err != nil {
			return err
		}
		if err := decoder(values[j], fv); err != nil {
			return err
		}
	}
	reflext.Indirect(v).Set(reflext.Indirect(vv))
	if !c.Next() {
		defer c.Close()
	}
	return nil
}

// All :
func (c *Cursor) All(results interface{}) error {
	defer c.Close()
	if c.err != nil {
		return c.err
	}

	v := reflext.ValueOf(results)
	if !reflext.IsKind(v.Type(), reflect.Ptr) {
		return ErrUnaddressableEntity
	}

	v = reflext.Indirect(v)
	t := v.Type()
	if !reflext.IsKind(t, reflect.Slice) {
		return xerrors.New("it must be a slice of entity")
	}

	length := len(c.columns)
	slice := reflect.MakeSlice(t, 0, 0)
	t = t.Elem()
	mapper := core.DefaultMapper
	idxs := mapper.TraversalsByName(t, c.columns)
	decoders := make([]codec.ValueDecoder, length, length)

	for i := 0; c.rows.Next(); i++ {
		values := make([]interface{}, length, length)
		for j := 0; j < length; j++ {
			values[j] = &values[j]
		}
		if err := c.rows.Scan(values...); err != nil {
			return err
		}

		vv := reflext.Zero(t)
		for j, idx := range idxs {
			fv := mapper.FieldByIndexes(vv, idx)
			if i < 1 {
				decoder, err := c.registry.LookupDecoder(fv.Type())
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
func (c *Cursor) Error() error {
	if c.rows != nil {
		defer c.rows.Close()
	}
	return c.err
}

// Next :
func (c *Cursor) Next() bool {
	return c.rows.Next()
}

// Close :
func (c *Cursor) Close() error {
	if c.rows != nil {
		defer c.rows.Close()
		return nil
	}
	return nil
}
