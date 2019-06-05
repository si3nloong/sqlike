package jsonb

// ReadNumber :
func (r *Reader) ReadNumber() (int64, error) {
	c := r.nextToken()
	if c == 'n' {
		r.unreadByte().ReadNull()
		return 0, nil
	}
	// if c == '-' || (c >= '0' && c <= '9') {
	// 	r.unreadByte()
	// 	for i := r.pos; i < r.len; i++ {
	// 		c = r.b[i]
	// 		if c >= '0' && c <= '9' {
	// 			continue
	// 		} else if c == '.' || c == 'e' {
	// 		} else {
	// 			str = string(r.b[r.pos:i])
	// 			r.pos = i
	// 			break
	// 		}
	// 	}
	// 	return str, nil
	// }
	return 0, nil
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
