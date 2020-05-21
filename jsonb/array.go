package jsonb

// ReadArray :
func (r *Reader) ReadArray(cb func(r *Reader) error) error {
	c := r.nextToken()
	if c == 'n' {
		return r.unreadByte().skipNull()
	}

	if c != '[' {
		return ErrInvalidJSON{
			callback: "ReadArray",
			message:  "expect start with [ for array",
		}
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

		if cb != nil {
			it := NewReader(b)
			if err := cb(it); err != nil {
				return err
			}
		}

		c = r.nextToken()
		if c != ',' {
			break
		}
	}

	if c != ']' {
		return ErrInvalidJSON{
			callback: "ReadArray",
			message:  "expect end with ] for array",
		}
	}
	return nil
}
