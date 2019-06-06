package jsonb 

import (
	"strings"
)


// ReadObject :
func (r *Reader) ReadObject(cb func(*Reader, string) error) error {
	c := r.nextToken()
	if c != '{' {
		return ErrDecode{}
	}

	var (
		k string
		err error
	)

	for {
		c = r.nextToken()
		if c == '}' {
			break
		}
		if c != '"' {
			panic("1")
		}
		k, err = r.unreadByte().ReadString()
		if err != nil {
			return err
		}
		c = r.nextToken()
		if c != ':' {
			panic("2")
		}
		// TODO: process the value
		if err := cb(r, k); err != nil {
			return err
		}
		c = r.nextToken()
		if c != ',' {
			break
		}
	}
	return nil
}

// ReadFlattenObject :
func (r *Reader) ReadFlattenObject(cb func(*Reader, string) error) error {
	level := 1
	c := r.nextToken()
	if c != '{' {
		return ErrDecode{}
	}

	var (
		paths []string
		key   string
		err error
	)

keyLoop:
	for {
		c = r.nextToken()
		if c == '}' {
			r.unreadByte()
			goto valueLoop
		}

		if c != '"' {
			return ErrDecode{}
		}
		key, err  = r.unreadByte().ReadString()
		if err != nil {
			return err
		}
		paths = append(paths, key)
		c = r.nextToken()
		if c != ':' {
			return ErrDecode{}
		}

		c = r.nextToken()
		switch c {
		case '{':
			level++
			goto keyLoop

		default:
			v, err := r.unreadByte().ReadBytes()
			if err != nil {
				return err
			}
			k := strings.Join(paths, ".")
			it := NewReader(v)
			if err := cb(it, k); err != nil {
				return err
			}
		}

	valueLoop:
		c = r.nextToken()
		switch c {
		case '}':
			level--
			if level < 1 {
				break keyLoop
			}
			paths = paths[:level-1]
			c = r.nextToken()
			if c != ',' {
				r.unreadByte()
			}

		case ',':
			paths = append([]string{}, paths[:len(paths)-1]...)

		default:
			break

		}
	}

	if c != '}' {
		return ErrDecode{}
	}
	return nil
}