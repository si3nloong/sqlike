package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestB2s(t *testing.T) {
	str := `hello world`
	require.Equal(t, str, B2s([]byte(str)))
}
