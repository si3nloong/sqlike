package jsonb

// ReadArray :
func (r *Reader) ReadArray(cb func(r *Reader) error) error {
	c := r.nextToken()
	if c == 'n' {
		return r.unreadByte().skipNull()
	}

	if c != '[' {
		return ErrDecode{}
	}

	c = r.nextToken()
	if c == ']' { // empty array
		return nil
	}

	r.unreadByte()

	for {
		b, err := r.ReadBytes()
		if err != nil {
			return err
		}

		it := NewReader(b)
		if err := cb(it); err != nil {
			return err
		}

		c = r.nextToken()
		if c != ',' {
			break
		}
	}

	if c != ']' {
		return ErrDecode{}
	}
	return nil
}
