package jsonb

// ReadArray :
func (r *Reader) ReadArray(cb func(r *Reader) error) error {
	level := 1
	c := r.nextToken()
	if c != '[' {
		return ErrDecode{}
	}

valueLoop:
	for {
		c = r.nextToken()
		if c == '[' {
			level++
			goto valueLoop
		}

		b, err := r.unreadByte().ReadBytes()
		if err != nil {
			return err
		}

		it := NewReader(b)
		if err := cb(it); err != nil {
			return err
		}

		c = r.nextToken()
		if c == ']' {
			level--

			if level <= 0 {
				break valueLoop
			}

			c = r.nextToken()
		}

		if c != ',' {
			break
		}
	}

	if c != ']' {
		return ErrDecode{}
	}
	return nil
}
