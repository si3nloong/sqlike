package core

import "unsafe"

// CalcBase64Length :
func CalcBase64Length(numBytes int) int {
	modulus := numBytes % 3
	if modulus == 0 {
		return ((numBytes / 3) * 4)
	}
	return ((numBytes / 3) * 4) + 4
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
