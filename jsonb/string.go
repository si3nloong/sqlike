package jsonb

import (
	"github.com/si3nloong/sqlike/util"
)

// ReadEscapeString :
func (r *Reader) ReadEscapeString() (string, error) {
	c := r.nextToken()
	if c == 'n' {
		if err := r.unreadByte().ReadNull(); err != nil {
			return "", ErrInvalidJSON{}
		}
		return "", nil
	}

	if c != '"' {
		return "", ErrInvalidJSON{}
	}

	blr := util.AcquireString()
	defer util.ReleaseString(blr)

	for i := r.pos; i < r.len; {
		c = r.b[i]
		if c == '"' {
			r.pos = i + 1
			break
		}

		if c == '\\' {
			switch r.b[i+1] {
			case '"':
				blr.WriteRune('"')
				i += 2
			case '\\':
				blr.WriteRune('\\')
				i += 2
			case 'b':
				blr.WriteRune('\b')
				i += 2
			case 'f':
				blr.WriteRune('\f')
				i += 2
			case 'n':
				blr.WriteRune('\n')
				i += 2
			case 'r':
				blr.WriteRune('\r')
				i += 2
			case 't':
				blr.WriteRune('\t')
				i += 2
			case '/':
				blr.WriteRune('/')
				i += 2
			case 'u':
				i += 2
			}
			continue
		}

		blr.WriteByte(c)
		i++
	}

	if c != '"' {
		return "", ErrInvalidJSON{}
	}

	return blr.String(), nil
}

// ReadString :
func (r *Reader) ReadString() (string, error) {
	c := r.nextToken()
	if c == 'n' {
		if err := r.unreadByte().ReadNull(); err != nil {
			return "", ErrInvalidJSON{}
		}
		return "", nil
	}

	if c != '"' {
		return "", ErrInvalidJSON{}
	}

	for i := r.pos; i < r.len; i++ {
		c = r.b[i]
		if c == '"' {
			str := string(r.b[r.pos:i])
			r.pos = i + 1
			return str, nil
		} else if c == '\\' {
			break
		} else if c < ' ' {
			// panic("unexpected character")
		}
	}

	return "", ErrInvalidJSON{}
}

// skipString :
func (r *Reader) skipString() error {
	c := r.nextToken()
	if c == 'n' {
		return r.unreadByte().ReadNull()
	}

	if c != '"' {
		return ErrInvalidJSON{}
	}

	for i := r.pos; i < r.len; {
		c = r.b[i]

		if c == '"' {
			r.pos = i + 1
			break
		}

		if c == '\\' {
			switch r.b[i+1] {
			case '"':
				i += 2
			case '\\':
				i += 2
			case 'b':
				i += 2
			case 'f':
				i += 2
			case 'n':
				i += 2
			case 'r':
				i += 2
			case 't':
				i += 2
			case '/':
				i += 2
			case 'u':
				i += 2
			}
			continue
		}
		i++
	}

	if c != '"' {
		return ErrInvalidJSON{}
	}

	return nil
}

var escapeCharMap = map[byte][]byte{
	'\t': []byte(`\t`),
	'\n': []byte(`\n`),
	'\r': []byte(`\r`),
	'"':  []byte(`\"`),
	'\\': []byte(`\\`),
	'\b': []byte(`\b`),
}

func escapeString(w *Writer, str string) {
	length := len(str)
	for i := 0; i < length; i++ {
		b := str[0]
		str = str[1:]
		if x, isOk := escapeCharMap[b]; isOk {
			w.Write(x)
			continue
		}
		w.WriteByte(b)
	}
}

func unescapeString(w *Writer, b string) {
	length := len(b)
	for i := 0; i < length; {
		c := b[i]

		if c == '\\' {
			switch b[i+1] {
			case '"':
				w.WriteRune('"')
				i += 2
			case '\\':
				w.WriteRune('\\')
				i += 2
			case 'b':
				w.WriteRune('\b')
				i += 2
			case 'f':
				w.WriteRune('\f')
				i += 2
			case 'n':
				w.WriteRune('\n')
				i += 2
			case 'r':
				w.WriteRune('\r')
				i += 2
			case 't':
				w.WriteRune('\t')
				i += 2
			case '/':
				w.WriteRune('/')
				i += 2
			case 'u':
				i += 2
			}
		} else {
			w.WriteByte(c)
			i++
		}
	}
}
