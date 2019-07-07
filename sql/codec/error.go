package codec

import (
	"reflect"
)

// ErrNoEncoder :
type ErrNoEncoder struct {
	Type reflect.Type
}

func (err ErrNoEncoder) Error() (msg string) {
	if err.Type == nil {
		msg = "no encoder for <nil>"
		return
	}
	msg = "no encoder for " + err.Type.String()
	return
}

// ErrNoDecoder :
type ErrNoDecoder struct {
	Type reflect.Type
}

func (err ErrNoDecoder) Error() (msg string) {
	if err.Type == nil {
		msg = "no decoder for <nil>"
		return
	}
	msg = "no decoder for " + err.Type.String()
	return
}
