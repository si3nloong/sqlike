package jsonb

import (
	"errors"
)

var whiteSpaceMap = map[byte]bool{
	' ':  true,
	'\n': true,
	'\t': true,
	'\r': true,
}

// Reader :
type Reader struct {
	typ   jsonType
	b     []byte
	pos   int
	len   int
	start int
	end   int
}

// NewReader :
func NewReader(b []byte) *Reader {
	length := len(b)
	copier := make([]byte, length)
	copy(copier, b)
	return &Reader{b: copier, len: length}
}

func (r *Reader) reset() *Reader {
	r.pos = 0
	return r
}

// Bytes :
func (r *Reader) Bytes() []byte {
	return r.b
}

// nextToken : nextToken will skip all whitespace and stop on a char
func (r *Reader) nextToken() byte {
	var c byte
	for i := r.pos; i < r.len; i++ {
		c = r.b[i]
		if _, ok := whiteSpaceMap[c]; ok {
			r.b = append(r.b[:i], r.b[i+1:]...)
			r.len = len(r.b)
			i--
			continue
		}
		r.pos = i + 1
		return c
	}
	return 0
}

func (r *Reader) prevToken() byte {
	if r.pos > 0 {
		return r.b[r.pos-1]
	}
	return 0
}

func (r *Reader) peekType() jsonType {
	c := r.nextToken()
	defer r.unreadByte()
	typ := valueMap[c]
	return typ
}

// IsNull :
func (r *Reader) IsNull() bool {
	offset := r.pos + 4
	if offset > r.len {
		return false
	}
	if string(r.b[r.pos:offset]) == null {
		return true
	}
	return false
}

func (r *Reader) skipArray() {
	level := 1
	c := r.nextToken()
	if c != '[' {
		return
	}

	for i := r.pos; i < r.len; i++ {
		switch r.b[i] {
		case '"': // If inside string, skip it
			// iter.head = i + 1
			// iter.skipString()
			// i = iter.head - 1 // it will be i++ soon
		case '[': // If open symbol, increase level
			level++
		case ']': // If close symbol, increase level
			level--

			// If we have returned to the original level, we're done
			if level <= 0 {
				r.pos = i + 1
				return
			}
		}
	}
}

// ReadBytes :
func (r *Reader) ReadBytes() ([]byte, error) {
	i := r.pos
	if err := r.skip(); err != nil {
		return nil, err
	}
	return r.b[i:r.pos], nil
}

// ReadValue :
func (r *Reader) ReadValue() (any, error) {
	typ := r.peekType()
	switch typ {
	case jsonString:
		return r.ReadString()
	case jsonNumber:
		num, err := r.ReadNumber()
		if err != nil {
			return nil, err
		}
		return num.Float64()
	case jsonBoolean:
		return r.ReadBoolean()
	case jsonNull:
		if err := r.ReadNull(); err != nil {
			return nil, err
		}
		return nil, nil
	case jsonArray:
		var v []any
		if err := r.ReadArray(func(it *Reader) error {
			x, err := it.ReadValue()
			if err != nil {
				return err
			}
			v = append(v, x)
			return nil
		}); err != nil {
			return v, err
		}
		return v, nil
	case jsonObject:
		var v map[string]any
		if err := r.ReadObject(func(it *Reader, k string) error {
			if v == nil {
				v = make(map[string]any)
			}
			x, err := it.ReadValue()
			if err != nil {
				return err
			}
			v[k] = x
			return nil
		}); err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, errors.New("invalid json format")
	}
}

func (r *Reader) unreadByte() *Reader {
	if r.pos > 0 {
		r.pos--
	}
	return r
}

// ReadBoolean :
func (r *Reader) ReadBoolean() (bool, error) {
	c := r.nextToken()
	r.unreadByte()
	if c == 'n' {
		if err := r.skipBytes([]byte{'n', 'u', 'l', 'l'}); err != nil {
			return false, err
		}
		return false, nil
	} else if c == 't' {
		if err := r.skipBytes([]byte{'t', 'r', 'u', 'e'}); err != nil {
			return false, err
		}
		return true, nil
	} else if c == 'f' {
		if err := r.skipBytes([]byte{'f', 'a', 'l', 's', 'e'}); err != nil {
			return false, err
		}
		return false, nil
	}
	return false, errors.New("invalid boolean value")
}

// ReadNull :
func (r *Reader) ReadNull() error {
	c := r.nextToken()
	if c == 'n' {
		return r.skipBytes([]byte{'u', 'l', 'l'})
	}
	return nil
}
