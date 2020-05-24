package sqldump

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDumper(t *testing.T) {
	dumper := NewDumper("driver", nil)
	require.NotNil(t, dumper)
}
