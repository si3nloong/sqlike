package mysql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type stringer struct {
}

func (s stringer) String() string {
	return "i'm stringer"
}

func TestFormat(t *testing.T) {
	var (
		ms  = New()
		str string
	)

	str = ms.Format(int64(-638731231286))
	require.Equal(t, "-638731231286", str)

	str = ms.Format(uint64(638731231286))
	require.Equal(t, "638731231286", str)

	str = ms.Format(int8(99))
	require.Equal(t, "99", str)

	str = ms.Format("hello world")
	require.Equal(t, `"hello world"`, str)

	str = ms.Format(`😆 😉 😊 😋 emojis`)
	require.Equal(t, `"😆 😉 😊 😋 emojis"`, str)

	str = ms.Format(true)
	require.Equal(t, "1", str)

	str = ms.Format(false)
	require.Equal(t, "0", str)

	str = ms.Format(float64(1232.888333))
	require.Equal(t, "1.232888333e+03", str)

	str = ms.Format(nil)
	require.Equal(t, "NULL", str)

	str = ms.Format(stringer{})
	require.Equal(t, `"i'm stringer"`, str)

	ts, _ := time.Parse("2006-01-02 15:04:05", "2020-01-03 12:00:40")
	str = ms.Format(ts)
	require.Equal(t, `"2020-01-03 12:00:40"`, str)
}
