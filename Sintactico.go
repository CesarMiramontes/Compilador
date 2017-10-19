// Sintactico
package main

import (
	"fmt"
	"strconv"
)

type NodeKind int

const (
	StmtK = 1 + iota
	ExpK
)

type StmtKind int

const (
	ProgramK = 1 + iota
	Ifk
	WriteK
	ReadK
	DeclK
	AssignK
	Dok
	UntilK
	WhileK
)

type ExpKind int

const (
	OpK = 1 + iota
	ConstK
	IdK
)

type ExpType int

const (
	Void = 1 + iota
	Integer
	Boolean
)

type Kind struct {
	stmt StmtKind
	exp  ExpKind
}

type Attr struct {
	op   token_types
	tipe string
	val  float64
	name string
}

type TreeNode struct {
	// guarda al token, sea operando u operador
	token Token
	// solo valido si el token es operador
	precedence int
	sibling    *TreeNode
	nodekind   NodeKind
	kind       Kind
	attr       Attr
	tipe       ExpType
	branch     [3]*TreeNode
}

func match(expected token_types) {
	if expected == token.tipo {
		token = getToken2()
	} else {
		isCorrect = false

		if token.tipo == TKN_ERROR {
			token = getToken2()
		}

		fmt.Println("sysntaxis error en match", token.lexema, token.linea, token.tipo, expected)
	}
}

func newStmtNode(kind StmtKind) *TreeNode {
	var t *TreeNode
	t = new(TreeNode)
	t.token = token

	for i := 0; i < 3; i++ {
		t.nodekind = StmtK
		t.kind.stmt = kind

	}

	return t
}

func newExpNode(expKind ExpKind) *TreeNode {
	var t *TreeNode
	t = new(TreeNode)
	t.token = token

	for i := 0; i < 3; i++ {
		t.nodekind = ExpK
		t.kind.exp = expKind

	}

	return t
}

func stmt_sequence() *TreeNode {
	t := statement()
	p := t

	for token.tipo != TKN_RBRACE {
		var q *TreeNode
		//match(TKN_SEMICOLON)

		q = statement()

		if q != nil {
			if t == nil {
				p = q
				t = p

			} else {
				p.sibling = q
				p = q

			}
		}
	}

	match(TKN_RBRACE)
	return t
}

func if_stmt_sequence() *TreeNode {
	t := statement()
	p := t
	flag := true

	if token.tipo == TKN_ELSE {
		flag = false
	} else if token.tipo == TKN_RBRACE {
		flag = false
		match(TKN_RBRACE)
	} else {

	}

	for flag {
		var q *TreeNode
		match(TKN_SEMICOLON)

		q = statement()

		if q != nil {
			if t == nil {
				p = q
				t = p

			} else {
				p.sibling = q
				p = q

			}
		}

		if token.tipo == TKN_ELSE {
			flag = false
		} else if token.tipo == TKN_RBRACE {
			flag = false
			match(TKN_RBRACE)
		} else {

		}

	}

	return t
}

func statement() *TreeNode {
	var t *TreeNode

	switch token.tipo {
	case TKN_ID:
		t = assign_stmt()
		match(TKN_SEMICOLON)

	case TKN_IF:
		t = if_stmt()

	//case TKN_ELSE:

	case TKN_READ:
		t = read_stmt()
		match(TKN_SEMICOLON)

	case TKN_WRITE:
		t = write_stmt()
		match(TKN_SEMICOLON)

	case TKN_TIPO:
		t = declaration_stmt()
		match(TKN_SEMICOLON)

	case TKN_RBRACE:
		//match(TKN_RBRACE)
		//fmt.Println(token)

	case TKN_WHILE:
		t = while_stmt()

	case TKN_DO:
		t = do_stmt()
		match(TKN_SEMICOLON)

	case TKN_PROGRAM:
		t = program_stmt()
		//match(TKN_RBRACE)

	default:
		fmt.Println("syntaxis error en stmt", token.lexema, token.tipo)
		token = getToken2()
	}

	return t
}

func program_stmt() *TreeNode {
	t := newStmtNode(ProgramK)
	match(TKN_PROGRAM)
	match(TKN_LBRACE)

	if t != nil {
		t.branch[0] = stmt_sequence()
	}

	return t
}

func if_stmt() *TreeNode {
	t := newStmtNode(Ifk)
	//recurda cambiar el if de match
	match(TKN_IF)

	if t != nil {
		match(TKN_LPAREN)
		t.branch[0] = expresion()
		match(TKN_RPAREN)
	}

	match(TKN_LBRACE)

	if t != nil {
		t.branch[1] = stmt_sequence()
	}

	//match(TKN_RBRACE)

	if token.lexema == "else" {
		match(TKN_ELSE)
		match(TKN_LBRACE)

		if t != nil {
			t.branch[2] = stmt_sequence()
		}
	}

	//match(TKN_RBRACE)
	match(TKN_FI)

	if token.tipo == TKN_RBRACE {

	} else {
		t.sibling = statement()
	}

	return t
}

func while_stmt() *TreeNode {
	t := newStmtNode(WhileK)
	//recurda cambiar el if de match
	match(TKN_WHILE)

	if t != nil {
		match(TKN_LPAREN)
		t.branch[0] = expresion()
		match(TKN_RPAREN)
	}

	match(TKN_LBRACE)

	if t != nil {
		t.branch[1] = stmt_sequence()
	}

	return t
}

func do_stmt() *TreeNode {
	t := newStmtNode(Dok)
	//recurda cambiar el if de match
	match(TKN_DO)
	match(TKN_LBRACE)

	if t != nil {
		t.branch[0] = stmt_sequence()
	}

	match(TKN_UNTIL)

	if t != nil {
		match(TKN_LPAREN)
		t.branch[1] = expresion()
		match(TKN_RPAREN)
	}

	return t
}

func declaration_stmt() *TreeNode {
	t := newStmtNode(DeclK)
	tipe := ""

	switch token.lexema {
	case "float":
		tipe = "float"
		break

	case "bool":
		tipe = "bool"
		break

	case "int":
		tipe = "int"
		break

	}

	//fmt.Println(tipe)

	match(TKN_TIPO)
	w := newStmtNode(AssignK)
	t.branch[0] = w
	flag := false

	for token.tipo == TKN_ID || token.tipo == TKN_COMMA {
		var q *TreeNode
		p := newStmtNode(AssignK)

		for e := t.branch[0]; e != nil; e = e.sibling {
			q = e
		}

		if token.tipo == TKN_ID {
			if !flag {
				q.attr.name = token.lexema
				q.attr.tipe = tipe
				//fmt.Println(p, q)
				//q = p
			} else {
				p.attr.name = token.lexema
				p.attr.tipe = tipe
				//fmt.Println(p,q)
				q.sibling = p
			}

		} else {
			fmt.Println("syntaxis error declarationstmt", token.lexema, token.linea)
		}

		/*if token.tipo == TKN_ID {
			if !flag {
				p.attr.name = token.lexema
				p.attr.tipe = tipe
				//fmt.Println(p,q)
				q = p
			} else {
				p.attr.name = token.lexema
				p.attr.tipe = tipe
				//fmt.Println(p,q)
				q.sibling = p
			}

		} else {
			fmt.Println("syntaxis error declarationstmt", token.lexema, token.linea)
		}*/

		match(TKN_ID)

		if token.tipo == TKN_COMMA {
			match(TKN_COMMA)
		}

		flag = true

	}

	//fmt.Println(t.branch[0].attr, t.branch[0].sibling.attr)
	//match(TKN_SEMICOLON)

	return t
}

func assign_stmt() *TreeNode {
	q := newStmtNode(AssignK)

	if q != nil && token.tipo == TKN_ID {
		q.attr.name = token.lexema
	}

	match(TKN_ID)
	t := newStmtNode(AssignK)

	if t != nil && token.tipo == TKN_ASSIGN {
		t.attr.name = token.lexema
	}

	match(TKN_ASSIGN)

	if t != nil {
		t.branch[0] = q
		t.branch[1] = expresion()
	}

	return t
}

func read_stmt() *TreeNode {
	t := newStmtNode(ReadK)
	match(TKN_READ)

	if t != nil && token.tipo == TKN_ID {
		t.attr.name = token.lexema
		t.branch[0] = newExpNode(IdK)
	}

	match(TKN_ID)
	return t
}

func write_stmt() *TreeNode {
	t := newStmtNode(WriteK)
	match(TKN_WRITE)

	if t != nil {
		t.branch[0] = expresion()
	}

	return t
}

func isLogic() bool {
	if token.tipo == TKN_LESS || token.tipo == TKN_ELESS {
		return true

	} else if token.tipo == TKN_EQUAL || token.tipo == TKN_NEQUAL {
		return true

	} else if token.tipo == TKN_MORE || token.tipo == TKN_EMORE {
		return true

	} else {
		return false
	}
}

func expresion() *TreeNode {
	t := exp()

	for token.tipo == TKN_AND || token.tipo == TKN_OR {
		p := newExpNode(OpK)

		if p != nil {
			p.branch[0] = t
			p.attr.op = token.tipo
			t = p
		}

		match(token.tipo)

		if t != nil {
			t.branch[1] = exp()
		}
	}

	return t
}

func exp() *TreeNode {
	t := simple_exp()

	if isLogic() {
		p := newExpNode(OpK)

		if p != nil {
			p.branch[0] = t
			p.attr.op = token.tipo
			t = p
		}

		match(token.tipo)

		if t != nil {
			t.branch[1] = simple_exp()
		}
	}

	return t
}

func simple_exp() *TreeNode {
	t := term()

	for token.tipo == TKN_ADD || token.tipo == TKN_MINUS {
		p := newExpNode(OpK)

		if p != nil {
			p.branch[0] = t
			p.attr.op = token.tipo
			t = p

			match(token.tipo)

			t.branch[1] = term()
		}
	}

	return t
}

func term() *TreeNode {
	t := factor()

	for token.tipo == TKN_MULTI || token.tipo == TKN_DIV {
		p := newExpNode(OpK)

		if p != nil {
			p.branch[0] = t
			p.attr.op = token.tipo
			t = p

			match(token.tipo)

			p.branch[1] = factor()
		}
	}

	return t
}

func factor() *TreeNode {
	var t *TreeNode

	switch token.tipo {
	case TKN_NUM:
		t = newExpNode(ConstK)

		if t == nil && token.tipo == TKN_NUM {
			variable, erro := strconv.ParseFloat(token.lexema, 64)

			if erro != nil {
				t.attr.val = variable
			}
		}

		match(TKN_NUM)

	case TKN_ID:
		t = newExpNode(IdK)

		if t == nil && token.tipo == TKN_ID {
			t.attr.name = token.lexema
		}

		match(TKN_ID)

	case TKN_ADD:
		t = newExpNode(TKN_MORE)

		if t == nil && token.tipo == TKN_ID {
			t.attr.name = token.lexema
		}

		match(TKN_MORE)
		t.branch[0] = factor()

	case TKN_MINUS:
		t = newExpNode(TKN_MINUS)

		if t == nil && token.tipo == TKN_ID {
			t.attr.name = token.lexema
		}

		match(TKN_MINUS)
		t.branch[0] = factor()

	case TKN_NOT:
		t = newExpNode(TKN_NOT)

		if t == nil && token.tipo == TKN_ID {
			t.attr.name = token.lexema
		}

		match(TKN_NOT)
		t.branch[0] = factor()

	case TKN_LPAREN:
		match(TKN_LPAREN)
		t = exp()
		match(TKN_RPAREN)

	default:
		fmt.Println("sysntaxis error en factor", token.lexema, token.linea)
		token = getToken2()

	}

	return t
}

func parse() *TreeNode {
	var t *TreeNode
	token = getToken2()
	t = stmt_sequence()
	//fmt.Println(t.token)
	//t = statement()
	//fmt.Println(auxVar)
	//token = getToken()
	/*t = stmt_sequence()
		if (token!=ENDFILE)
	    	syntaxError("Code ends before file\n");*/
	return t
}
