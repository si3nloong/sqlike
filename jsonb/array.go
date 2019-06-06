package jsonb

// ReadArray :
func (r *Reader) ReadArray(cb func(r *Reader) error) error {
	c := r.nextToken()
	if c != '[' {
		return ErrDecode{}
	}

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
