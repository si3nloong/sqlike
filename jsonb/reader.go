package jsonb

import (
	"log"

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
	valueMap['\n'] = jsonWhitespace
	valueMap['['] = jsonArray
	valueMap['{'] = jsonObject
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
	ast   *Token
	b     []byte
	Value *Token
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

func (r *Reader) setNode(c byte) {
	if r.ast == nil { // root
		r.ast = new(Token)
		r.ast.b = []byte{c}
		log.Println("Root initialize", *r.ast)
		return
	}
	n := new(Token)
	r.ast.child = n
}

// ReadNext :
func (r *Reader) nextToken() byte {
	var c byte
	for i := r.pos; i < r.len; i++ {
		c = r.b[i]
		if _, isOk := whiteSpaceMap[c]; isOk {
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

// ReadString :
func (r *Reader) ReadString() (str string) {
	c := r.nextToken()
	if c != '"' {
		panic("it should be string")
	}
	if c == 'n' {
		r.unreadByte()
		r.ReadNull()
	}
	// r.pos++
	for i := r.pos; i < r.len; i++ {
		c = r.b[i]
		if c == '"' {
			str = string(r.b[r.pos:i])
			r.pos = i + 1
			return
		} else if c == '\\' {
			break
		} else if c < ' ' {
			panic("unexpected character")
		}
	}
	return
}

// ReadNumber :
func (r *Reader) ReadNumber() (str string) {
	c := r.nextToken()
	if c == '-' || (c >= '0' && c <= '9') {
		r.unreadByte()
		for i := r.pos; i < r.len; i++ {
			c = r.b[i]
			if c >= '0' && c <= '9' {
				continue
			} else if c == '.' || c == 'e' {
			} else {
				str = string(r.b[r.pos:i])
				r.pos = i
				break
			}
		}
		return
	}
	if c == 'n' {
		r.unreadByte()
		r.ReadNull()
	}
	return
}

// ReadBoolean :
func (r *Reader) ReadBoolean() (bool, error) {
	c := r.nextToken()
	log.Println("Booolean", string(c))
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
		return ErrUnexpectedChar{}
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
		return ErrUnexpectedChar{}
	}

	var (
		path string
		key  string
	)

keyLoop:
	for {
		c = r.nextToken()
		if c != '"' {
			return ErrUnexpectedChar{}
		}
		key = r.unreadByte().ReadString()
		if path != "" {
			key = path + "." + key
		}
		c = r.nextToken()
		if c != ':' {
			return ErrUnexpectedChar{}
		}

		c = r.nextToken()
		switch c {
		case '{':
			level++
			path = key
			continue keyLoop

		default:
			v, err := r.unreadByte().ReadBytes()
			if err != nil {
				return err
			}
			m[key] = v
		}

		c = r.nextToken()
		if c == ',' {
			continue keyLoop
		}
		for level > 1 {
			if c != '}' {
				return ErrUnexpectedChar{}
			}
			path = ""
			c = r.nextToken()
			level--
		}

		if c != ',' {
			break
		}
	}

	if c != '}' {
		return ErrUnexpectedChar{}
	}
	log.Println("Debug ============>")
	for k, v := range m {
		rdr := NewReader(v)
		if err := cb(rdr, k); err != nil {
			return err
		}
	}
	return nil
}
