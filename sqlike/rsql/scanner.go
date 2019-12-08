package rsql

import (
	"errors"
	"log"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/timtadh/lexmachine"
)

// Scanner :
type Scanner struct {
	level  int
	token  *lexmachine.Token
	parser *Parser
	*lexmachine.Scanner
}

// NextToken :
func (scan *Scanner) NextToken() (*lexmachine.Token, bool) {
	it, err, eof := scan.Next()
	for !eof && err == nil && it == nil {
		it, err, eof = scan.Next()
	}
	if err != nil || eof {
		return nil, true
	}
	return it.(*lexmachine.Token), eof
}

// ParseTokens :
func (scan *Scanner) ParseTokens() (primitive.Group, error) {
	grp := new(primitive.Group)
	grps := make([]bool, 0)

	for {
		tkn, eof := scan.NextToken()
		if eof {
			break
		}

		char := string(tkn.Lexeme)
		switch tkn.Type {
		case Whitespace:
			// skip
		case Group:
			if char == "(" {
				grps = append(grps, true)
				grp.Values = append(grp.Values, expr.Raw("("))
			} else {
				grps = grps[:len(grps)-1]
				grp.Values = append(grp.Values, expr.Raw(")"))
			}
		case Text:
			// check expression
			// scan.ParseExpression(grp, tkn)
			log.Println("Text :", char)
		case String:
			// scan.ParseExpression(grp, tkn)
			log.Println("String :", char)
		case And:
			grp.Values = append(grp.Values, primitive.And)
		case Or:
			grp.Values = append(grp.Values, primitive.Or)
		}
	}

	log.Println("Filters :::", grp.Values)
	// tkn, eof := scan.NextToken()
	// if eof {
	// 	return nil
	// }
	// // scan.token = tkn
	// char := string(tkn.Lexeme)
	// if tkn.Type == Group && char == "(" {
	// 	scan.level++
	// 	scan.values.Values = append(scan.values.Values, expr.Raw(char))
	// 	return scan.ParseGroup()
	// }
	// return scan.ParseExpression(tkn)
	return *grp, nil
}

// ParseExpression :
func (scan *Scanner) ParseExpression(grp *primitive.Group, column *lexmachine.Token) error {
	log.Println("Mapper :", scan.parser.mapper.Names)
	log.Println("Key :", string(column.Lexeme))
	field, ok := scan.parser.mapper.Names[string(column.Lexeme)]
	if !ok {
		return errors.New("")
	}

	if _, ok := field.Tag.LookUp("filter"); !ok {
		return errors.New("")
	}

	operator, eof := scan.NextToken()
	if eof {
		return errors.New("")
	}

	value, eof := scan.NextToken()
	if eof {
		return errors.New("")
	}

	decoder, err := scan.parser.registry.LookupDecoder(field.Type)
	if err != nil {
		return err
	}

	v := reflext.Zero(field.Type)
	log.Println("datatype :", v.Type())
	if err := decoder(string(value.Lexeme), v); err != nil {
		log.Println("Decode err:", err)
		return err
	}

	log.Println("token 2 :", string(operator.Lexeme))
	log.Println("value :", value)

	it := v.Interface()
	log.Println("after :", it)

	// (==|!=|>|>=|<|<=|=ne=|=nin=)
	name := field.Name
	if v, ok := field.Tag.LookUp("column"); ok {
		name = v
	}
	switch string(operator.Lexeme) {
	case "==":
		grp.Values = append(grp.Values, expr.Equal(name, it))
	case "!=":
		grp.Values = append(grp.Values, expr.NotEqual(name, it))
	case ">":
		grp.Values = append(grp.Values, expr.GreaterThan(name, it))
	case ">=":
		grp.Values = append(grp.Values, expr.GreaterOrEqual(name, it))
	case "<":
		grp.Values = append(grp.Values, expr.LesserThan(name, it))
	case "<=":
		grp.Values = append(grp.Values, expr.LesserOrEqual(name, it))
	case "=ne=":
		grp.Values = append(grp.Values, expr.NotIn(name, it))
	case "=nin=":
		grp.Values = append(grp.Values, expr.NotIn(name, it))
	}

	return nil
}

// var (
// 	valueMap = make([]tokenType, 256)
// )

// // EOF :
// var EOF = byte(0)

// // Token :
// type Token struct {
// 	t tokenType
// 	s string
// }

// // Scanner :
// type Scanner struct {
// 	r *bufio.Reader
// }

// // NewScanner :
// func NewScanner() *Scanner {
// 	return &Scanner{}
// }

// // Scan returns the next token and literal value.
// func (s *Scanner) Scan(r io.Reader) (out []Token, err error) {
// 	s.r = bufio.NewReader(r)
// 	// log.Println(s.r.ReadRune())

// 	for {
// 		tok, lit := s.ScanToken()
// 		log.Println("Value ::", tok, lit)
// 		if tok == EOF {
// 			break
// 		}
// 		// else if tok == ILLEGAL {
// 		// 	return out, fmt.Errorf("Illegal Token : %s", lit)
// 		// } else {
// 		// 	out = append(out, NewTokenString(tok, lit))
// 	}
// 	// }

// 	return
// }

// func (s *Scanner) read() byte {
// 	b, err := s.r.ReadByte()
// 	if err != nil {
// 		return EOF
// 	}
// 	return b
// }

// // ScanToken :
// func (s *Scanner) ScanToken() (tok byte, lit string) {
// 	b := s.read()

// 	log.Println("Byte :", string(b), b)
// 	switch b {
// 	case ' ', '\n', '\t', '\r':
// 		// skip
// 	default:
// 		tok = b
// 	}

// 	// if isReservedRune(ch) {
// 	// 	s.unread()
// 	// 	return s.scanReservedRune()
// 	// } else if isIdent(ch) {
// 	// 	s.unread()
// 	// 	return s.scanIdent()
// 	// }

// 	// if ch == eof {
// 	// 	return EOF, ""
// 	// }

// 	// return ILLEGAL, string(ch)
// 	return
// }

// func (s *Scanner) scanTill() {

// }

// func (s *Scanner) scanString() {
// 	for {

// 	}
// }

// func init() {
// 	length := len(valueMap)
// 	for i := 0; i < length; i++ {
// 		valueMap[i] = Invalid
// 	}
// 	valueMap['"'] = String
// 	valueMap['-'] = Numeric
// 	valueMap['0'] = Numeric
// 	valueMap['1'] = Numeric
// 	valueMap['2'] = Numeric
// 	valueMap['3'] = Numeric
// 	valueMap['4'] = Numeric
// 	valueMap['5'] = Numeric
// 	valueMap['6'] = Numeric
// 	valueMap['7'] = Numeric
// 	valueMap['8'] = Numeric
// 	valueMap['9'] = Numeric
// 	valueMap[' '] = Whitespace
// 	valueMap['\r'] = Whitespace
// 	valueMap['\t'] = Whitespace
// 	valueMap['\n'] = Whitespace

// 	s := NewScanner()
// 	log.Println("Scanner :", s)
// 	s.Scan(bytes.NewReader([]byte(`testing==%asdasjdkl`)))
// }
