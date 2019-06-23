package sqlike

import "golang.org/x/xerrors"

// errors :
var (
	ErrNoValueUpdate = xerrors.New("sqlike: no value to update")
	// ErrInvalidInput :
	ErrInvalidInput = xerrors.New("sqlike: invalid input <nil>")
	// ErrUnaddressableEntity :
	ErrUnaddressableEntity = xerrors.New("sqlike: unaddressable entity")
	// ErrNilEntity :
	ErrNilEntity = xerrors.New("sqlike: entity is <nil>")
)
