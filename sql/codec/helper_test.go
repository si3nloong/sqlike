package codec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestB2s(t *testing.T) {
	str := `hello world`
	require.Equal(t, str, b2s([]byte(str)))
}
