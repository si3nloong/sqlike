package jsonb

import (
	"bytes"
	"log"
)

var emptyJSON = []byte(`null`)

// Reader :
type Reader struct {
	typ    jsonType
	b      []byte
	length int
}

// ReadNext :
func (r *Reader) ReadNext() {
	var b byte
	for i := 0; i < r.length; {
		b = r.b[i]
		if b == ' ' {
			i++
			continue
		}

		if b == '\\' {
			switch r.b[i+1] {
			case 't':
				i += 2
			case 'r':
				i += 2
			case 'n':
				i += 2
			}
		}

		switch b {
		case '[':
			r.typ = jsonArray
		case '{':
			r.typ = jsonObject
		case 'n':
			if bytes.Equal(r.b[i:i+4], emptyJSON) {
			}
		case 't':
			if bytes.Equal(r.b[i:i+4], []byte(`true`)) {
			}
		case 'f':
			if bytes.Equal(r.b[i:i+5], []byte(`false`)) {
			}
		case '"':
			r.typ = jsonString
		}
		log.Println("debug :", string(b))
		i++
	}
}
