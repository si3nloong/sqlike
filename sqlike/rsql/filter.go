package rsql

import "log"

func (p *Parser) parseFilter(values map[string]string, params *Params) (errs Errors) {
	val, ok := values[p.FilterTag]
	log.Println("QueryString :", string(val))
	if !ok || len(val) < 1 {
		return nil
	}

	lxr, _ := p.lexer.Scanner([]byte(val))
	scan := &Scanner{parser: p, Scanner: lxr}
	if err := scan.ParseToken(); err != nil {
		log.Println("Error :", err)
	}
	params.Filters = scan.values
	return nil
}
