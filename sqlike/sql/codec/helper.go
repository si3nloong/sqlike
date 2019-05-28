package codec

import "unsafe"

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
