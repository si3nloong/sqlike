package jsonb

import (
	"bytes"

	"golang.org/x/xerrors"
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
		return xerrors.New("unexpected json length")
	}
	if !bytes.Equal(b, r.b[r.pos:offset]) {
		return ErrInvalidJSON{}
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
	return
}

func (r *Reader) skipObject() {
	level := 1
	c := r.nextToken()
	if c != '{' {
		return
	}

	// TODO : index out of range
	for level > 0 {
		c = r.nextToken()
		switch c {
		case '"':
		case '{':
			level++
		case '}':
			level--
		}
	}

	// for i := r.pos; i < r.len; i++ {
	// 	switch r.b[i] {
	// 	case '"': // If inside string, skip it
	// 		// iter.head = i + 1
	// 		r.pos = i
	// 		r.skipString()
	// 		i = r.pos + 1
	// 		// i = iter.head - 1 // it will be i++ soon
	// 	case '{': // If open symbol, increase level
	// 		level++
	// 	case '}': // If close symbol, increase level
	// 		level--
	// 		log.Println("Pos", i)

	// 		// If we have returned to the original level, we're done
	// 		if level == 0 {
	// 			r.pos = i + 1
	// 			log.Println(r.pos, r.peekType())
	// 			return
	// 		}
	// 	}
	// }
}
