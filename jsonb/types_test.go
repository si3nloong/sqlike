package jsonb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsonType(t *testing.T) {

	require.Equal(t, "invalid", jsonInvalid.String())
	require.Equal(t, "null", jsonNull.String())
	require.Equal(t, "object", jsonObject.String())
	require.Equal(t, "array", jsonArray.String())
	require.Equal(t, "whitespace", jsonWhitespace.String())
	require.Equal(t, "string", jsonString.String())
	require.Equal(t, "boolean", jsonBoolean.String())
	require.Equal(t, "number", jsonNumber.String())

}
