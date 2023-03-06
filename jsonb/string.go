package jsonb

import (
	"github.com/si3nloong/sqlike/v2/internal/util"
)

// ReadRawString :
func (r *Reader) ReadRawString() (string, error) {
	c := r.nextToken()
	if c == 'n' {
		if err := r.unreadByte().ReadNull(); err != nil {
			return "", err
		}
		return "null", nil
	}

	if c != '"' {
		return "", ErrInvalidJSON{
			callback: "ReadString",
			message:  "expect start with \"",
		}
	}

	for i := r.pos; i < r.len; i++ {
		c = r.b[i]
		if c == '"' {
			str := string(r.b[r.pos:i])
			r.pos = i + 1
			return str, nil
		}
	}

	return "", ErrInvalidJSON{
		callback: "ReadString",
		message:  "expect end with \"",
	}
}

// ReadString :
func (r *Reader) ReadString() (string, error) {
	c := r.nextToken()
	if c == 'n' {
		if err := r.unreadByte().ReadNull(); err != nil {
			return "", err
		}
		return "", nil
	}

	if c != '"' {
		return "", ErrInvalidJSON{
			callback: "ReadString",
			message:  "expect start with \"",
		}
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
				// TODO: unicode
				blr.WriteRune('u')
				i += 2
			default:
				blr.WriteByte(c)
			}
			continue
		}

		blr.WriteByte(c)
		i++
	}

	if c != '"' {
		return "", ErrInvalidJSON{
			callback: "ReadEscapeString",
			message:  "expect end with \"",
		}
	}

	return blr.String(), nil
}

// skipString :
func (r *Reader) skipString() error {
	c := r.nextToken()
	if c == 'n' {
		return r.unreadByte().ReadNull()
	}

	if c != '"' {
		return ErrInvalidJSON{
			callback: "skipString",
			message:  "expect start with \"",
		}
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
		return ErrInvalidJSON{
			callback: "skipString",
			message:  "expect end with \"",
		}
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

func escapeString(w JsonWriter, str string) {
	length := len(str)
	for i := 0; i < length; i++ {
		b := str[0]
		str = str[1:]
		if x, ok := escapeCharMap[b]; ok {
			w.Write(x)
			continue
		}
		w.WriteByte(b)
	}
}
