package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	rawConn := "user:password@/dbname"
	opt := Connect()
	opt.ApplyURI(rawConn)

	{
		require.Equal(t, rawConn, opt.raw)
		require.Equal(t, rawConn, opt.RawConnStr())
	}

	{
		opt.SetUsername("root")
		require.Equal(t, "root", opt.Username)
	}

	{
		opt.SetPassword("abcd1234")
		require.Equal(t, "abcd1234", opt.Password)
	}

	{
		opt.SetProtocol("tcp")
		require.Equal(t, "tcp", opt.Protocol)
	}

	{
		opt.SetHost("127.0.0.1")
		require.Equal(t, "127.0.0.1", opt.Host)

		opt.SetHost("192.168.0.10")
		require.Equal(t, "192.168.0.10", opt.Host)
	}

	{
		opt.SetPort("3307")
		require.Equal(t, "3307", opt.Port)

		require.Panics(t, func() {
			opt.SetPort("127.0.0.1")
		})

		require.Panics(t, func() {
			opt.SetPort("3306 7707")
		})
	}

	{
		opt.SetSocket("unix()")
		require.Equal(t, "unix()", opt.Socket)
	}

	{
		opt.SetCharset("utf8")
		require.Equal(t, "utf8", string(opt.Charset))
	}

	{
		opt.SetCollate("utf8_bin_general")
		require.Equal(t, "utf8_bin_general", opt.Collate)
	}
}
