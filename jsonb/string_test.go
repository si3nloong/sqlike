package jsonb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadEscapeString(t *testing.T) {

}

func TestReadRawString(t *testing.T) {
	t.Run("ReadRawString with invalid null", func(t *testing.T) {
		r := NewReader([]byte(`nan`))
		str, err := r.ReadRawString()
		require.Error(t, err)
		require.Empty(t, str)
	})

	t.Run("ReadRawString with null", func(t *testing.T) {
		r := NewReader([]byte(`null`))
		str, err := r.ReadRawString()
		require.NoError(t, err)
		require.Equal(t, `null`, str)
	})

	t.Run("ReadRawString", func(t *testing.T) {
		r := NewReader([]byte(`"Hello World !"`))
		str, err := r.ReadRawString()
		require.NoError(t, err)
		require.Equal(t, `Hello World !`, str)
	})

	t.Run("ReadRawString with unquoted string", func(t *testing.T) {
		r := NewReader([]byte(`"12344`))
		str, err := r.ReadRawString()
		require.Error(t, err)
		require.Empty(t, str)
	})
}

func TestReadString(t *testing.T) {

	t.Run("escapeString", func(t *testing.T) {
		str := `test		\"asdasjkd
	lkasd128378127#$%^&*()_)(*&^%$#@#~!@#$%`
		blr := NewWriter()
		escapeString(blr, str)
		assert.Equal(t, `test\t\t\\\"asdasjkd\n\tlkasd128378127#$%^&*()_)(*&^%$#@#~!@#$%`, blr.String(), "unexpected result")
	})
	t.Run("ReadString with escape string", func(t *testing.T) {
		r := NewReader([]byte(`"1234\"4abc\"de\nfg"`))
		str, err := r.ReadString()
		require.NoError(t, err)
		require.Equal(t, "1234\"4abc\"de\nfg", str)
	})

	t.Run("ReadString with unquoted string", func(t *testing.T) {
		r := NewReader([]byte(`"12344`))
		str, err := r.ReadString()
		require.Error(t, err)
		require.Empty(t, str)
	})

	t.Run("ReadString with null", func(t *testing.T) {
		r := NewReader([]byte(`null`))
		str, err := r.ReadString()
		require.NoError(t, err)
		require.Empty(t, str)
	})

	t.Run("ReadString with errors", func(t *testing.T) {
		r := NewReader([]byte(`ss""`))
		str, err := r.ReadString()
		require.Error(t, err)
		require.Empty(t, str)
	})
}
