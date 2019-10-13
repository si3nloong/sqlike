package rsql

import (
	"errors"
	"log"

	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
	"github.com/timtadh/lexmachine"
)

// Scanner :
type Scanner struct {
	values primitive.Group
	level  int
	token  *lexmachine.Token
	parser *Parser
	*lexmachine.Scanner
}

// NextToken :
func (scan *Scanner) NextToken() (*lexmachine.Token, bool) {
	it, _, eof := scan.Next()
	for !eof && it == nil {
		it, _, eof = scan.Next()
	}
	// log.Println("Token :", it)
	return it.(*lexmachine.Token), eof
}

// ParseToken :
func (scan *Scanner) ParseToken() error {
	tkn, eof := scan.NextToken()
	if eof {
		return nil
	}
	// scan.token = tkn
	char := string(tkn.Lexeme)
	if tkn.Type == Group && char == "(" {
		scan.level++
		scan.values.Values = append(scan.values.Values, expr.Raw(char))
		return scan.ParseGroup()
	}
	return scan.ParseExpression(tkn)
}

// ParseExpression :
func (scan *Scanner) ParseExpression(tkn *lexmachine.Token) error {
	expression, eof := scan.NextToken()
	if eof {
		return errors.New("")
	}

	value, eof := scan.NextToken()
	if eof {
		return errors.New("")
	}
	// log.Println("Debug Expression ===========>")

	log.Println(string(tkn.Lexeme), string(expression.Lexeme), string(value.Lexeme))
	tkn3, eof := scan.NextToken()
	if eof {
		return nil
	}
	scan.values.Values = append(scan.values.Values, expr.Equal(string(tkn.Lexeme), string(value.Lexeme)))
	log.Println(tkn3.Type == Group)
	switch tkn3.Type {
	case Group:
		scan.TC = tkn3.TC
	case And:
		scan.values.Values = append(scan.values.Values, primitive.And)
	case Or:
		scan.values.Values = append(scan.values.Values, primitive.Or)
	default:
		return errors.New("invalid syntax")
	}
	return nil
}

// ParseGroup :
func (scan *Scanner) ParseGroup() error {
	for {
		log.Println("HERE !!!")
		tkn, eof := scan.NextToken()
		if eof {
			break
		}
		char := string(tkn.Lexeme)
		if tkn.Type == Group && char == ")" {
			scan.level--
			scan.values.Values = append(scan.values.Values, expr.Raw(char))
			break
		}
		if err := scan.ParseExpression(tkn); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	// b := []byte(`(_id==133,(category!=-10.00;num==.922;test=="value\""));d1=="";c1==testing,d1!=108)`)
	// lexer := lexmachine.NewLexer()
	// // skip white space
	// lexer.Add([]byte(`\s`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
	// 	return nil, nil
	// })
	// lexer.Add([]byte(`\(|\)`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
	// 	return &Token{Type: Group, Lexeme: match.Bytes}, nil
	// })
	// lexer.Add([]byte(`\"(\\.|[^\"])*\"`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
	// 	return &Token{Type: String, Lexeme: match.Bytes}, nil
	// })
	// lexer.Add([]byte(`(\,|or)`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
	// 	return &Token{Type: Or, Lexeme: match.Bytes}, nil
	// })
	// lexer.Add([]byte(`(\;|and)`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
	// 	return &Token{Type: And, Lexeme: match.Bytes}, nil
	// })
	// lexer.Add([]byte(`(\-)?([0-9]*\.[0-9]+|[0-9]+)`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
	// 	return &Token{Type: Numeric, Lexeme: match.Bytes}, nil
	// })
	// lexer.Add([]byte(`[a-zA-Z0-9\_\.\%]+`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
	// 	return &Token{Type: Text, Lexeme: match.Bytes}, nil
	// })
	// lexer.Add([]byte(`(\=\=|\!\=|\>|\>\=|\<|\<\=|\=ne\=|\=nin\=)`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
	// 	return &Token{Type: Operator, Lexeme: match.Bytes}, nil
	// })

	// s, _ := lexer.Scanner(b)
	// group := make([]*Token, 0)
	// values := make([]interface{}, 0)
	// log.Println(lexer)
	// for {
	// 	tkn, err, eof := s.Next()
	// 	if err != nil {

	// 	}
	// 	if eof {
	// 		log.Println(eof, err)
	// 		break
	// 	}

	// 	if tkn == nil {
	// 		continue
	// 	}
	// 	field := tkn.(*Token)

	// 	if field.Type == Group {
	// 		group = append(group, field)
	// 		continue
	// 	}

	// 	tkn, _, _ = s.Next()
	// 	operator := tkn.(*Token)
	// 	if operator.Type != Operator {
	// 		break
	// 	}

	// 	tkn, _, _ = s.Next()
	// 	value = tkn.(*Token)

	// 	switch string(operator.Lexeme) {
	// 	case "!=":
	// 	}

	// 	values = append(values)

	// 	tkn, _, _ = s.Next()
	// 	x = tkn.(*Token)
	// 	if x.Type != And && x.Type != Or && x.Type != Group {
	// 		break
	// 	}
	// 	log.Println(x.Type, "|", string(x.Lexeme), x)

	// 	// log.Println("Token :", token)
	// }
	// log.Println()
	// log.Println(s.TC)
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
