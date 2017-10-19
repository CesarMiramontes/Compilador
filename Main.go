// Main
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

var auxVar int

var tokens []Token

var token Token

var isCorrect bool

var primitivo string

var tableSymbols = make(map[int]map[string]int)

var symbolKind = make(map[int]map[string]string)

var wf *bufio.Writer

func writeTree(path string, root *TreeNode) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)

	mostrarArbol(root, w)

	/*for _, line := range lines {
		fmt.Fprintln(w, line)
	}*/

	return w.Flush()
}

func getToken2() Token {
	var tok Token

	if auxVar < len(tokens) {
		if tokens[auxVar].lexema == "" {
			if auxVar+1 < len(tokens) {
				auxVar++
				tok = tokens[auxVar]
			} else {
				tok = Token{TKN_RBRACE, "", 0, 0}
			}
		} else {
			tok = tokens[auxVar]
		}

		tok = tokens[auxVar]

	} else {
		tok = Token{TKN_RBRACE, "", 0, 0}
		return tok
	}

	auxVar++
	return tok
}

func toArray(lista *list.List) []Token {
	var tokenArray []Token

	for e := lista.Front(); e != nil; e = e.Next() {
		tok := e.Value.(Token)
		tokenArray = append(tokenArray, tok)
	}

	return tokenArray
}

func mostrarArbol(root *TreeNode, w *bufio.Writer) {

	for e := root.branch[0]; e != nil; e = e.sibling {
		//fmt.Println(e.token.lexema)
		fmt.Fprintln(w, e.token.lexema)
		mostrarHijos(e, "", w)
	}

}

func mostrarHijos(root *TreeNode, tabulacion string, w *bufio.Writer) {
	for i := 0; i < 3; i++ {
		if root.branch[i] != nil {
			mostrarHermanos(root.branch[i], tabulacion+"	", w)
		}
	}
}

func mostrarHermanos(root *TreeNode, tabulacion string, w *bufio.Writer) {
	for e := root; e != nil; e = e.sibling {
		//fmt.Println(tabulacion, e.token.lexema)
		fmt.Fprintln(w, tabulacion, e.token.lexema)
		mostrarHijos(e, tabulacion+"	", w)
	}
}

func main() {
	auxVar = 0
	isCorrect = true
	lista := list.New()

	//Llamamos al analizador lexico
	//y nos tegresa una lista de tokens
	Lexico(" ", lista)
	tokens = toArray(lista)

	if isCorrect {
		fmt.Println("Exito: el analisis Lexico fue exitoso :)")
	} else {
		panic("Error en Lexico, no se realizara analsis Sintactico X_X")

	}

	//Llamamos a el analizador sintactico
	//consecuencia de esto nos regresara la raiz de un arbol fmt.Println(t.sibling.branch[1].branch[1].branch[0].token)
	t := parse()

	writeTree("tree.txt", t)

	if isCorrect {
		fmt.Println("Exito: el analisis Sintactico fue exitoso :)")
	} else {
		panic("Error en Sintactico, no se realizara analsis Semantico X_X")

	}

	semantico(t)

	if isCorrect {
		fmt.Println("Exito: el analisis Semantico fue exitoso :)")
	} else {
		panic("Error en Semantico, no se generara codigo intermedio X_X")
	}

	code_generator(t)

}
