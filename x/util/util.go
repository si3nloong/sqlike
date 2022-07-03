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
	return *(*string)(unsafe.Pointer(&b))
}

// BytesToString converts byte slice to string.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to byte slice.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
