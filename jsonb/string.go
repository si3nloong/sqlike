package jsonb

// ReadString :
func (r *Reader) ReadString() (string, error) {
	c := r.nextToken()
	if c == 'n' {
		r.unreadByte().ReadNull()
		return "", nil
	}

	if c != '"' {
		return "", ErrDecode{}
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
	return "", ErrDecode{}
}

// skipString :
func (r *Reader) skipString() {
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
