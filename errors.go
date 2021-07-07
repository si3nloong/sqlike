package sqlike

import "errors"

// errors : common error of sqlike
var (
	ErrNoValueUpdate = errors.New("sqlike: no value to update")
	// ErrInvalidInput :
	ErrInvalidInput = errors.New("sqlike: invalid input <nil>")
	// ErrUnaddressableEntity :
	ErrUnaddressableEntity = errors.New("sqlike: unaddressable entity")
	// ErrNilEntity :
	ErrNilEntity = errors.New("sqlike: entity is <nil>")
	// ErrNoColumn :
	ErrNoColumn = errors.New("sqlike: no columns to create index")
	// ErrNestedTransaction :
	ErrNestedTransaction = errors.New("sqlike: nested transaction")
)
