package jsonb

// ReadNumber :
func (r *Reader) ReadNumber() (string, error) {
	c := r.nextToken()
	if c == 'n' {
		r.unreadByte().ReadNull()
		return "0", nil
	}

	if valueMap[c] != jsonNumber {
		return "", ErrInvalidJSON{}
	}

	r.unreadByte()
	str := string(r.b[r.pos:])
	for i := r.pos; i < r.len; i++ {
		c = r.nextToken()
		if c != '.' && c != 'e' && valueMap[c] != jsonNumber {
			str = string(r.b[r.pos:i])
			r.pos = i
			break
		}
	}

	return str, nil
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
	return
}
