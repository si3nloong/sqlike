package jsonb

import (
	"bytes"
	"fmt"

	"errors"
)

func (r *Reader) skip() (err error) {
	typ := r.peekType()
	switch typ {
	case jsonString:
		err = r.skipString()
	case jsonNumber:
		r.skipNumber()
	case jsonBoolean:
		r.skipBoolean()
	case jsonNull:
		err = r.skipBytes([]byte{'n', 'u', 'l', 'l'})
	case jsonArray:
		r.skipArray()
	case jsonObject:
		r.skipObject()
	}
	return err
}

func (r *Reader) skipBytes(b []byte) error {
	offset := r.pos + len(b)
	if offset > r.len {
		return errors.New("unexpected json length")
	}
	cutset := r.b[r.pos:offset]
	if !bytes.Equal(b, cutset) {
		return ErrInvalidJSON{
			callback: "skipBytes",
			message:  fmt.Sprintf("expected %s", b),
		}
	}
	r.pos = offset
	return nil
}

func (r *Reader) skipNull() error {
	c := r.nextToken()
	r.unreadByte()
	if c == 'n' {
		return r.skipBytes([]byte{'n', 'u', 'l', 'l'})
	}
	return ErrInvalidJSON{}
}

func (r *Reader) skipBoolean() {
	c := r.nextToken()
	r.unreadByte()
	if c == 'n' {
		r.skipBytes([]byte{'n', 'u', 'l', 'l'})
		return
	}
	if c == 't' {
		r.skipBytes([]byte{'t', 'r', 'u', 'e'})
		return
	}
	if c == 'f' {
		r.skipBytes([]byte{'f', 'a', 'l', 's', 'e'})
		return
	}
}

func (r *Reader) skipObject() error {
	c := r.nextToken()
	if c != '{' {
		return errors.New("object should start with {")
	}

loop:
	for {
		c = r.nextToken()

		// key: value
		switch c {
		case '}':
			break loop

		case '"': // expect key
			r.unreadByte()
			if err := r.skipString(); err != nil {
				return errors.New("object key must be string")
			}

			c = r.nextToken()
			if c != ':' {
				return errors.New("character : must place within key and value")
			}

			// skip anything
			if err := r.skip(); err != nil {
				return err
			}

			c = r.nextToken()
			switch c {
			case '}':
				break loop
			case ',':
				continue
			default:
				return fmt.Errorf("expected , or } after object value")
			}

		default:
			return fmt.Errorf("invalid character %s in object", string(c))
		}
	}

	if r.b[r.pos] != '}' {
		return errors.New("invalid char on end of object")
	}

	return nil
}
