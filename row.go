package sqlike

import (
	"database/sql"
)

// Row is the result of calling QueryRow to select a single row.
type Row struct {
	*sql.Row
}

// Decode will unmarshal the values to struct.
func (r *Row) Decode(dst any) error {
	// v := reflext.ValueOf(dst)
	// if !v.IsValid() {
	// 	return ErrInvalidInput
	// }

	// t := v.Type()
	// if !reflext.IsKind(t, reflect.Ptr) {
	// 	return ErrUnaddressableEntity
	// }

	// t = reflext.Deref(t)
	// if !reflext.IsKind(t, reflect.Struct) {
	// 	return errors.New("sqlike: it must be a struct to decode")
	// }

	// idxs := r.cache.TraversalsByName(t, r.columns)
	// values, err := r.values()
	// if err != nil {
	// 	return err
	// }
	// vv := reflext.Zero(t)
	// for j, idx := range idxs {
	// 	if idx == nil {
	// 		continue
	// 	}
	// 	fv := r.cache.FieldByIndexes(vv, idx)
	// 	decoder, err := r.dialect.LookupDecoder(fv.Type())
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if err := decoder(values[j], fv); err != nil {
	// 		return err
	// 	}
	// }
	// reflext.IndirectInit(v).Set(reflext.Indirect(vv))
	// if r.close {
	// 	return r.Close()
	// }
	return nil
}
