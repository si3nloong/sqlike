package jsonb

import (
	"encoding/json"
	"fmt"
)

// Number :
type Number = json.Number

// ReadNumber :
func (r *Reader) ReadNumber() (Number, error) {
	c := r.nextToken()
	if c == 'n' {
		if err := r.unreadByte().ReadNull(); err != nil {
			return "", err
		}
		return "0", nil
	}

	if valueMap[c] != jsonNumber {
		return "", ErrInvalidJSON{
			callback: "ReadNumber",
			message:  fmt.Sprintf("invalid character %q, expected number", c),
		}
	}

	r.unreadByte()
	str := string(r.b[r.pos:])
	pos := r.pos
	for i := pos; i < r.len; i++ {
		c = r.nextToken()
		if c != '.' && c != 'e' && valueMap[c] != jsonNumber {
			str = string(r.b[pos:i])
			r.pos = i
			break
		}
	}

	return json.Number(str), nil
}

func (r *Reader) skipNumber() {
	c := r.nextToken()
	if c == 'n' {
		r.unreadByte().ReadNull()
		return
	}

	for i := r.pos; i < r.len; i++ {
		c = r.b[i]
		switch c {
		case ' ', '\n', '\r', '\t', ',', '}', ']':
			r.pos = i
			return
		}
	}
}
