package sqldump

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDumper(t *testing.T) {
	dumper := NewDumper("mysql", nil)
	require.NotNil(t, dumper)
}

func TestBackupTo(t *testing.T) {

}

func TestGetVersion(t *testing.T) {
	// dumper := NewDumper("driver", nil)
	// v, err := dumper.getVersion(context.Background())
	// require.NoError(t, err)
	// require.NotEmpty(t, v)
}

func TestGetColumns(t *testing.T) {

}
