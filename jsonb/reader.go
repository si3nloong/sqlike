package jsonb

import (
	"golang.org/x/xerrors"
)

var whiteSpaceMap = map[byte]bool{
	' ':  true,
	'\n': true,
	'\t': true,
	'\r': true,
}

var emptyJSON = []byte(`null`)

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
	copier := make([]byte, length, length)
	copy(copier, b)
	return &Reader{b: copier, len: length}
}

// Bytes :
func (r *Reader) Bytes() []byte {
	return r.b
}

// ReadNext :
func (r *Reader) nextToken() byte {
	var c byte
	for i := r.pos; i < r.len; i++ {
		c = r.b[i]
		if _, isOk := whiteSpaceMap[c]; isOk {
			r.b = append(r.b[:i], r.b[i+1:]...)
			r.len = r.len - 1
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

// GetBytes :
func (r *Reader) GetBytes() (b []byte) {
	r.start = r.pos
	c := r.nextToken()
	switch c {
	case '"':
		// r.skipString()
	case 'n':
		// r.skipThreeBytes('u', 'l', 'l') // null
	case 't':
		// iter.skipThreeBytes('r', 'u', 'e') // true
	case 'f':
		// iter.skipFourBytes('a', 'l', 's', 'e') // false
	case '0':
		// iter.unreadByte()
		// iter.ReadFloat32()
	case '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// iter.skipNumber()
	case '[':
		// r.skipArray()
	case '{':
		// r.skipObject()
	default:
		// iter.ReportError("Skip", fmt.Sprintf("do not know how to skip: %v", c))
		return
	}
	return
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
func (r *Reader) ReadValue() (interface{}, error) {
	typ := r.peekType()
	switch typ {
	case jsonString:
		return r.ReadEscapeString()
	case jsonNumber:
		return r.ReadNumber()
	case jsonBoolean:
		return r.ReadBoolean()
	case jsonNull:
		if err := r.ReadNull(); err != nil {
			return nil, err
		}
		return nil, nil
	case jsonArray:
		var v []interface{}
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
		var v map[string]interface{}
		if err := r.ReadObject(func(it *Reader, k string) error {
			if v == nil {
				v = make(map[string]interface{})
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
		return nil, xerrors.New("invalid json format")
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
		r.skipBytes([]byte{'n', 'u', 'l', 'l'})
		return false, nil
	}
	if c == 't' {
		r.skipBytes([]byte{'t', 'r', 'u', 'e'})
		return true, nil
	}
	if c == 'f' {
		r.skipBytes([]byte{'f', 'a', 'l', 's', 'e'})
		return false, nil
	}
	return false, xerrors.New("invalid boolean value")
}

// ReadNull :
func (r *Reader) ReadNull() error {
	c := r.nextToken()
	if c == 'n' {
		return r.skipBytes([]byte{'u', 'l', 'l'})
	}
	return nil
}
