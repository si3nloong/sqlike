package jsonb

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	{
		err := ErrInvalidJSON{callback: "funcName", message: "error"}
		require.Equal(t, "jsonb.funcName: error", err.Error())
	}

	{
		err := ErrNoEncoder{}
		require.Equal(t, "no encoder for <nil>", err.Error())
	}

	{
		err := ErrNoEncoder{Type: reflect.TypeOf("")}
		require.Equal(t, "no encoder for string", err.Error())
	}

	{
		err := ErrNoDecoder{}
		require.Equal(t, "no decoder for <nil>", err.Error())
	}

	{
		err := ErrNoDecoder{Type: reflect.TypeOf(int(10))}
		require.Equal(t, "no decoder for int", err.Error())
	}
}
