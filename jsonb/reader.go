package jsonb

import (
	"strings"

	"golang.org/x/xerrors"
)

var whiteSpaceMap = map[byte]bool{
	' ':  true,
	'\n': true,
	'\t': true,
	'\r': true,
}

var valueMap = make([]jsonType, 256)

func init() {
	valueMap['"'] = jsonString
	valueMap['-'] = jsonNumber
	valueMap['0'] = jsonNumber
	valueMap['1'] = jsonNumber
	valueMap['2'] = jsonNumber
	valueMap['3'] = jsonNumber
	valueMap['4'] = jsonNumber
	valueMap['5'] = jsonNumber
	valueMap['6'] = jsonNumber
	valueMap['7'] = jsonNumber
	valueMap['8'] = jsonNumber
	valueMap['9'] = jsonNumber
	valueMap['t'] = jsonBoolean
	valueMap['f'] = jsonBoolean
	valueMap['n'] = jsonNull
	valueMap['['] = jsonArray
	valueMap['{'] = jsonObject
	valueMap[' '] = jsonWhitespace
	valueMap['\r'] = jsonWhitespace
	valueMap['\t'] = jsonWhitespace
	valueMap['\n'] = jsonWhitespace
}

var emptyJSON = []byte(`null`)

// Token :
type Token struct {
	typ   jsonType
	b     []byte
	child *Token
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
	return &Reader{b: b, len: len(b)}
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

func (r *Reader) skipArray() {
	level := 1
	for {
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
				if level == 0 {
					r.pos = i + 1
					return
				}
			}
		}
	}
}

// ReadBytes :
func (r *Reader) ReadBytes() ([]byte, error) {
	i := r.pos
	r.skip()
	return r.b[i:r.pos], nil
}

// ReadValue :
func (r *Reader) ReadValue() (interface{}, error) {
	typ := r.peekType()
	switch typ {
	case jsonString:
		return r.ReadString(), nil
	case jsonNumber:
		return r.ReadNumber(), nil
	case jsonBoolean:
		return r.ReadBoolean()
	case jsonNull:
		return r.ReadNull(), nil
	case jsonArray:
		var v []interface{}
		return v, nil
	case jsonObject:
		var v map[string]interface{}
		if err := r.ReadObject(func(r *Reader, k string) error {
			if v == nil {
				v = make(map[string]interface{})
			}
			x, err := r.ReadValue()
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
		return nil, nil
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
	if c == 't' {
		r.skipBytes([]byte{'r', 'u', 'e'})
		return true, nil
	}
	if c == 'f' {
		r.skipBytes([]byte{'a', 'l', 's', 'e'})
		return false, nil
	}
	return false, xerrors.New("invalid boolean value")
}

// ReadNull :
func (r *Reader) ReadNull() (b bool) {
	c := r.nextToken()
	if c == 'n' {
		r.skipBytes([]byte{'u', 'l', 'l'})
		return true
	}
	return
}

// ReadArray :
func (r *Reader) ReadArray() (arr []interface{}) {
	c := r.nextToken()
	for c == ',' {
		c = r.nextToken()
	}
	return
}

// ReadObject :
func (r *Reader) ReadObject(cb func(*Reader, string) error) error {
	c := r.nextToken()
	if c != '{' {
		return ErrDecode{}
	}

	var k string
	for {
		c = r.nextToken()
		if c == '}' {
			break
		}
		if c != '"' {
			panic("1")
		}
		k = r.unreadByte().ReadString()
		c = r.nextToken()
		if c != ':' {
			panic("2")
		}
		// TODO: process the value
		if err := cb(r, k); err != nil {
			return err
		}
		c = r.nextToken()
		if c != ',' {
			break
		}
	}
	return nil
}

// ReadFlattenObject :
func (r *Reader) ReadFlattenObject(cb func(*Reader, string) error) error {
	level := 1
	m := make(map[string][]byte)
	c := r.nextToken()
	if c != '{' {
		return ErrDecode{}
	}

	var (
		paths []string
		key   string
	)

keyLoop:
	for {
		c = r.nextToken()
		if c == '}' {
			r.unreadByte()
			goto valueLoop
		}

		if c != '"' {
			return ErrDecode{}
		}
		key = r.unreadByte().ReadString()
		paths = append(paths, key)
		c = r.nextToken()
		if c != ':' {
			return ErrDecode{}
		}

		c = r.nextToken()
		switch c {
		case '{':
			level++
			goto keyLoop

		default:
			v, err := r.unreadByte().ReadBytes()
			if err != nil {
				return err
			}
			k := strings.Join(paths, ".")
			m[k] = v
		}

	valueLoop:
		c = r.nextToken()
		switch c {
		case '}':
			level--
			if level < 1 {
				break keyLoop
			}
			paths = paths[:level-1]
			c = r.nextToken()
			if c != ',' {
				r.unreadByte()
			}

		case ',':
			paths = append([]string{}, paths[:len(paths)-1]...)

		default:
			break

		}
	}

	if c != '}' {
		return ErrDecode{}
	}

	for k, v := range m {
		// log.Println(k, ":", string(v))
		rdr := NewReader(v)
		if err := cb(rdr, k); err != nil {
			return err
		}
	}
	return nil
}
