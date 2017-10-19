// Codigo Intermedio
package main

import (
	"fmt"
	"strconv"
)

var t_index int

var loop_index int

var current_t string

var aux_t string

func code_assambler(t *TreeNode) {
	if t == nil {
		return
	}

	auxiliar_loop := loop_index

	switch t.nodekind {
	case StmtK:
		//fmt.Println("Stmtk:", t.token.lexema)

		switch t.kind.stmt {

		case ProgramK:
			//fmt.Println("program")
			fmt.Println("Main:")
			code_assambler(t.branch[0])
			break

		case Ifk:
			//fmt.Println("	IfZ", viewTree(t.branch[0]), "Goto loop"+strconv.Itoa(auxiliar_loop))
			fmt.Println("(if_f, " + viewTree(t.branch[0]) + ", " + strconv.Itoa(auxiliar_loop) + ", _)")
			loop_index++
			aux_if := loop_index
			//fmt.Println()
			code_assambler(t.branch[1])
			//fmt.Println("	Goto loop"+strconv.Itoa(aux_if), ":")
			fmt.Println("(_, _, loop" + strconv.Itoa(aux_if) + ", _)")
			//fmt.Println("loop"+strconv.Itoa(auxiliar_loop), ":")
			fmt.Println("(lab, " + strconv.Itoa(auxiliar_loop) + ", _, _)")
			code_assambler(t.branch[2])
			//fmt.Println("loop"+strconv.Itoa(aux_if), ":")
			fmt.Println("(lab, loop" + strconv.Itoa(aux_if) + ", _, _)")
			//kill_instance_varibles(deep + 1)
			break

		case WriteK:
			//fmt.Println("write")
			viewTree(t.branch[0])
			fmt.Println("(wr, " + strconv.Itoa(t_index) + ", _, _)")
			//fmt.Println("	PushParam __t" + strconv.Itoa(t_index))
			//fmt.Println("	LCall _PrintFloat")
			//fmt.Println("	PopParam")
			//fmt.Println(t.branch[0].token.lexema, ":=", viewTree(t.branch[1]))
			//fmt.Println()
			//token = t.branch[0].token

			//validate_exp_tree(t.branch[0], deep)
			//fmt.Println()
			break

		case ReadK:
			//fmt.Println("read")
			//	t_index++
			//	rr := "__t" + strconv.Itoa(t_index)
			//fmt.Println("	" + rr + " := LCall _ReadVariable")
			//fmt.Println("	" + rr + " := " + t.branch[0].token.lexema)
			fmt.Println("(rd, " + t.branch[0].token.lexema + ", _, _)")
			break

		case AssignK:
			//fmt.Println("assign")
			//fmt.Println("	"+t.branch[0].token.lexema, ":=", viewTree(t.branch[1]))
			fmt.Println("(asn, " + t.branch[0].token.lexema + ", _, " + viewTree(t.branch[1]) + ")")
			//fmt.Println()
			break

		case DeclK:
			//fmt.Println("declaration")
			//makeVariable(t.branch[0], deep)
			break

		case Dok:
			//fmt.Println("loop:", loop_index)
			fmt.Println("(lab, loop" + strconv.Itoa(loop_index) + "_, _)")
			loop_index++
			code_assambler(t.branch[0])
			//fmt.Println("	IfZ", viewTree(t.branch[1]), "Goto loop"+strconv.Itoa(auxiliar_loop))
			fmt.Println("(if_f, " + viewTree(t.branch[1]) + ", loop" + strconv.Itoa(auxiliar_loop) + ", _)")
			//fmt.Println()
			//kill_instance_varibles(deep + 1)
			break

		case UntilK:
			//fmt.Println("until")
			break

		case WhileK:
			//fmt.Println("while")
			viewTree(t.branch[0])
			//fmt.Println()
			code_assambler(t.branch[1])
			//kill_instance_varibles(deep + 1)
			break

		default:

		}
		break

	case ExpK:
		//fmt.Println("Expk:", t.token.lexema)

		break
	}

	code_assambler(t.sibling)

}

func viewTree(t *TreeNode) string {
	if t == nil {
		//fmt.Println()
		return ""
	}

	switch t.kind.exp {
	case OpK:
		r1 := viewTree(t.branch[0])
		r2 := viewTree(t.branch[1])
		t_index++
		r3 := "__t" + strconv.Itoa(t_index)
		//fmt.Println("	"+r3, ":=", r1, t.token.lexema, r2)
		make_4(t, r1, r2, r3)
		return r3
		break

	case IdK:
		return t.token.lexema

	case ConstK:
		return t.token.lexema

		break
	}

	return ""

	//viewTree(t.branch[0])
	//viewTree(t.branch[1])
	//fmt.Print(t.token.lexema, " ")
}

func make_4(t *TreeNode, r1 string, r2 string, r3 string) {
	if t.token.lexema == "+" {
		fmt.Println("(add, " + r1 + ", " + r2 + ", " + r3 + ")")

	} else if t.token.lexema == "-" {
		fmt.Println("(sub, " + r1 + ", " + r2 + ", " + r3 + ")")

	} else if t.token.lexema == "*" {
		fmt.Println("(mul, " + r1 + ", " + r2 + ", " + r3 + ")")

	} else if t.token.lexema == "/" {
		fmt.Println("(div, " + r1 + ", " + r2 + ", " + r3 + ")")

	} else {
		fmt.Println("(cmp, " + r1 + ", " + r2 + ", " + r3 + ")")

	}
}

func code_generator(t *TreeNode) {
	fmt.Println("generando codigo...")
	t_index = 0
	loop_index = 0
	current_t = "" + strconv.Itoa(123)
	code_assambler(t)
}
