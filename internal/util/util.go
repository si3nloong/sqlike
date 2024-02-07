package util

import (
	"strings"
	"sync"
	"unsafe"
)

var (
	strBldrPool = &sync.Pool{
		New: func() any {
			return new(strings.Builder)
		},
	}
)

// AcquireString :
func AcquireString() *strings.Builder {
	return strBldrPool.Get().(*strings.Builder)
}

// ReleaseString :
func ReleaseString(x *strings.Builder) {
	if x != nil {
		defer strBldrPool.Put(x)
		x.Reset()
	}
}

// UnsafeString :
func UnsafeString(b []byte) string {
	// #nosec
	return *(*string)(unsafe.Pointer(&b))
}

// BytesToString converts byte slice to string.
func BytesToString(b []byte) string {
	// #nosec
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to byte slice.
func StringToBytes(s string) []byte {
	// #nosec
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func EscapeString(s string, escChar byte) string {
	buf := make([]byte, 0, 3*len(s)/2)
	f := func() {
		buf = append(buf, []byte{escChar, escChar}...)
	}
	for w := 0; len(s) > 0; s = s[w:] {
		r := rune(s[0])
		w = 1
		if r == rune(escChar) {
			f()
			continue
		}
		buf = append(buf, s[0])
	}
	return string(buf)
}
