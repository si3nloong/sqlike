package mysql

import (
	"testing"

	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

// TestConnect :
func TestConnect(t *testing.T) {
	var (
		ms  = MySQL{}
		str string
	)

	str = ms.Connect(&options.ConnectOptions{
		Username: "root",
		Host:     "localhost",
		Port:     "3306",
	})
	require.Equal(t, `root:@tcp(localhost:3306)/?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci`, str)

	uri := `root:@unix(localhost:3306)/?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci`
	opt := new(options.ConnectOptions)
	str = ms.Connect(opt.ApplyURI(uri))
	require.Equal(t, uri, str)
	require.Panics(t, func() {
		ms.Connect(nil)
	})

	require.Panics(t, func() {
		opt := new(options.ConnectOptions)
		ms.Connect(opt)
	})
}
