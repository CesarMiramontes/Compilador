// Semantico
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func makeVariable(t *TreeNode, deep int) {
	//fmt.Println(deep)
	if t == nil {
		return
	}

	_, ok := tableSymbols[deep]
	if ok {
		//existe
	} else {
		//no existe
		tableSymbols[deep] = make(map[string]int)
		symbolKind[deep] = make(map[string]string)
	}

	token := t.token
	i := tableSymbols[deep][token.lexema]
	if i == 0 {
		tableSymbols[deep][token.lexema] = 1
		symbolKind[deep][token.lexema] = t.attr.tipe
		writeTable("", ""+t.attr.tipe+" "+token.lexema)
		//fmt.Println(t.token)
	} else {
		fmt.Println("error, variable", token.lexema, "is already defined")
		isCorrect = false
	}

	makeVariable(t.sibling, deep)
}

func kill_instance_varibles(deep int) {
	_, ok := tableSymbols[deep]
	if ok {
		//existe
		mm := tableSymbols[deep]
		for key, value := range mm {
			//fmt.Println("Key:", key, "Value:", value)
			writeTable("", "Eliminada variable de instancia "+key)
			value = 1 + value
		}
		//fmt.Println()
		delete(tableSymbols, deep)
		delete(symbolKind, deep)
	} else {
		//no existe
	}
}

func does_variable_exist(t *TreeNode, deep int) bool {
	token := t.token
	_, ok := tableSymbols[deep]

	if ok {
		i := tableSymbols[deep][token.lexema]
		//existe
		if i == 0 {
			if deep == 0 {
				fmt.Println("Error: variable", token.lexema, "is not defined")
				isCorrect = false
				return false
			} else {
				return does_variable_exist(t, (deep - 1))
			}

		} else {
			return true

		}
	} else {
		//no existe
		if deep == 0 {
			fmt.Println("Error: variable", token.lexema, "is not defined")
			isCorrect = false
			return false
		} else {
			return does_variable_exist(t, (deep - 1))
		}
	}
}

func get_primitive(t *TreeNode, deep int) string {
	token := t.token
	_, ok := tableSymbols[deep]
	i := tableSymbols[deep][token.lexema]
	j := symbolKind[deep][token.lexema]

	if ok {

		//existe
		if i == 0 {
			if deep == 0 {
				fmt.Println("Error: variable", token.lexema, "is not defined")
				isCorrect = false
				return j
			} else {
				return get_primitive(t, (deep - 1))
			}

		} else {
			return j

		}
	} else {
		//no existe
		if deep == 0 {
			fmt.Println("Error: variable", token.lexema, "is not defined")
			isCorrect = false
			return j
		} else {
			return get_primitive(t, (deep - 1))
		}
	}
}

func is_type_correct(t *TreeNode, tipo int) bool {
	//fmt.Println(get_primitive(t, tipo), primitivo)

	switch primitivo {
	case "int":
		if get_primitive(t, tipo) == "int" {
			return true
		} else if get_primitive(t, tipo) == "float" {
			//perdidaPrecision = true
			return true
		} else {
			return false
		}
		break

	case "float":
		if get_primitive(t, tipo) == "int" {
			return true
		} else if get_primitive(t, tipo) == "float" {
			return true
		} else {
			return false
		}
		break

	case "bool":
		break
	}

	if get_primitive(t, tipo) == primitivo {
		return true
	} else {
		return false
	}
}

func is_operator_correct(tipo string) bool {
	if tipo == "bool" {
		return isLogicSemantic()
	} else {
		return !isLogicSemantic()
	}
}

func validate_exp_tree(t *TreeNode, deep int) {
	if t == nil {
		return
	}

	//fmt.Println(t.token.lexema, t.kind.exp)

	switch t.kind.exp {
	case OpK:
		//fmt.Println("operador")
		token = t.token

		if token.lexema == "/" && primitivo == "int" {
			fmt.Println("posible perdida de precision en linea ", token.lexema)
		}

		if !is_operator_correct(primitivo) {
			fmt.Println("el token", t.token.lexema, "no corresponde a", primitivo)
			isCorrect = false
			return
		}
		break

	case ConstK:
		//fmt.Println("constante")
		validate_consts(t.token)
		break

	case IdK:
		//fmt.Println("identificador")
		does_variable_exist(t, deep)
		if !is_type_correct(t, deep) {
			fmt.Println("Token equivocado en la asignacion", primitivo, t.token.lexema, "linea:", t.token.linea)
			isCorrect = false
		}
		break

	default:
		//fmt.Println(t.token.lexema, t.kind.exp)
	}

	validate_exp_tree(t.branch[0], deep)
	//fmt.Println(t.token.lexema, t.kind.exp)
	validate_exp_tree(t.branch[1], deep)
}

/*func validate_exp_tree(t *TreeNode, deep int) {
	if t == nil {
		return
	}

	validate_exp_tree(t.branch[0], deep)
	//fmt.Println(t.token.lexema, t.kind.exp)
	validate_exp_tree(t.branch[1], deep)
}*/

func validate_consts(t Token) bool {
	//fmt.Println(t)
	switch t.tipo {
	case TKN_NUM:
		if strings.Contains(t.lexema, ".") && primitivo == "int" {
			fmt.Println("Posible perdad de precision en linea ", t.linea)
		} else {
			//fmt.Printf("")
		}

		if primitivo == "bool" {
			if t.lexema == "1" {

			} else if t.lexema == "0" {

			} else {
				isCorrect = false
				return false
			}
		}

		return true
		break

	default:
		isCorrect = false
		return false
	}

	isCorrect = false
	return false
}

func node_secuence(t *TreeNode, deep int) {
	if t == nil {
		return
	}

	switch t.nodekind {
	case StmtK:
		//fmt.Println("Stmtk:", t.token.lexema)

		switch t.kind.stmt {

		case ProgramK:
			//fmt.Println("program")
			node_secuence(t.branch[0], (deep + 1))
			kill_instance_varibles(deep + 1)
			break

		case Ifk:
			//fmt.Println("if")
			validate_boolean_expresion(t.branch[0], (deep + 1))
			kill_instance_varibles(deep + 1)
			node_secuence(t.branch[1], (deep + 1))
			kill_instance_varibles(deep + 1)
			node_secuence(t.branch[2], (deep + 1))
			kill_instance_varibles(deep + 1)
			break

		case WriteK:
			//fmt.Println("write")
			validate_write_expresion(t.branch[0], deep)
			//token = t.branch[0].token

			//validate_exp_tree(t.branch[0], deep)
			//fmt.Println()
			break

		case ReadK:
			//fmt.Println("read")
			does_variable_exist(t.branch[0], deep)
			break

		case AssignK:
			//fmt.Println("assign")
			if does_variable_exist(t.branch[0], deep) {
				//fmt.Println(t.branch[0].token, t.branch[0].attr)
				primitivo = get_primitive(t.branch[0], deep)
				//fmt.Println(primitivo)
				if primitivo == "bool" {
					validate_boolean_expresion(t.branch[1], deep)
				} else {
					validate_exp_tree(t.branch[1], deep)
				}

				//fmt.Println()
			}
			break

		case DeclK:
			//fmt.Println("declaration")
			makeVariable(t.branch[0], deep)
			break

		case Dok:
			//fmt.Println("do")
			node_secuence(t.branch[0], (deep + 1))
			kill_instance_varibles(deep + 1)
			validate_boolean_expresion(t.branch[1], (deep + 1))
			kill_instance_varibles(deep + 1)
			break

		case UntilK:
			//fmt.Println("until")
			break

		case WhileK:
			//fmt.Println("while")
			validate_boolean_expresion(t.branch[0], (deep + 1))
			node_secuence(t.branch[1], (deep + 1))
			kill_instance_varibles(deep + 1)
			break

		default:

		}
		break

	case ExpK:
		//fmt.Println("Expk:", t.token.lexema)

		break
	}

	node_secuence(t.sibling, deep)

}

func isLogicSemantic() bool {
	if token.tipo == TKN_AND || token.tipo == TKN_OR {
		return true
	} else if token.tipo == TKN_LESS || token.tipo == TKN_ELESS {
		return true

	} else if token.tipo == TKN_EQUAL || token.tipo == TKN_NEQUAL {
		return true

	} else if token.tipo == TKN_MORE || token.tipo == TKN_EMORE {
		return true

	} else {
		return false
	}
}

func isLogicFirstOrder() bool {
	if token.tipo == TKN_AND || token.tipo == TKN_OR {
		return true

	} else {
		return false
	}
}

func isLogicSecondOrder() bool {
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

func validate_boolean_expresion(t *TreeNode, deep int) {
	token = t.token
	if isLogicFirstOrder() {
		validate_boolean_expresion(t.branch[0], deep)
		validate_boolean_expresion(t.branch[1], deep)

	} else if isLogicSecondOrder() {
		primitivo = "float"
		validate_exp_tree(t.branch[0], deep)
		validate_exp_tree(t.branch[1], deep)

	} else if t.kind.exp == IdK {
		primitivo = "bool"
		if does_variable_exist(t, deep) && get_primitive(t, deep) == "bool" {

		} else {
			fmt.Println("asignacion incoreecta en linea ", t.token.linea)
			isCorrect = false
		}

	} else if t.kind.exp == ConstK {
		if validate_consts(t.token) {

		} else {
			fmt.Println("asignacion incorrecta en linea ", t.token.linea)
			isCorrect = false
		}
	} else {
		fmt.Println("Uso incorrecto de expresion booleana en linea ", t.token.linea)
		isCorrect = false
	}
}

func validate_write_expresion(t *TreeNode, deep int) {
	token = t.token
	if isLogicFirstOrder() {
		validate_boolean_expresion(t.branch[0], deep)
		validate_boolean_expresion(t.branch[1], deep)

	} else if isLogicSecondOrder() {
		primitivo = "float"
		validate_exp_tree(t.branch[0], deep)
		validate_exp_tree(t.branch[1], deep)

	} else if t.kind.exp == IdK {
		//primitivo = "bool"
		if does_variable_exist(t, deep) {

		} else {
			fmt.Println("no existe la variable de la linea ", t.token.linea)
			isCorrect = false
		}

	} else if t.kind.exp == ConstK {

	} else if t.kind.exp == OpK {
		primitivo = "float"
		validate_exp_tree(t, deep)

	} else {
		fmt.Println("Uso incorrecto de expresion write en linea ", t.token.linea)
		isCorrect = false
	}
}

func writeTable(path string, message string) {
	if message == "" {
		file, err := os.Create(path)
		//defer file.Close()
		wf = bufio.NewWriter(file)

		if err != nil {
			//return err
		}

		return
	}

	fmt.Fprintln(wf, message)

	//return nil
}

func semantico(t *TreeNode) error {
	writeTable("variable.txt", "")
	node_secuence(t, 0)
	return wf.Flush()
	//fmt.Println("Hello World!")
}
