// Lexico
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"unicode/utf8"
)

//definimos una enumeracion para los tipos de tokens
type token_types int

const (
	TKN_BEGIN     = 1 + iota
	TKN_RESERVED  //2
	TKN_END       //3
	TKN_READ      //4
	TKN_WRITE     //5
	TKN_TIPO      //6
	TKN_IF        //7
	TKN_ELSE      //8
	TKN_FI        //9
	TKN_PROGRAM   //10
	TKN_ID        //11
	TKN_NUM       //12
	TKN_LPAREN    //13
	TKN_RPAREN    //14
	TKN_LBRACE    //15
	TKN_RBRACE    //16
	TKN_SEMICOLON //17
	TKN_COMMA     //18
	TKN_ASSIGN    //19
	TKN_LESS      //20
	TKN_ELESS     //21
	TKN_MORE      //22
	TKN_EMORE     //23
	TKN_EQUAL     //24
	TKN_NEQUAL    //25
	TKN_ADD       //26
	TKN_MINUS     //27
	TKN_MULTI     //28
	TKN_DIV       //29
	TKN_POW       //30
	TKN_EOF       //31
	TKN_ERROR     //32
	TKN_DO        //33
	TKN_UNTIL     //34
	TKN_WHILE     //35
	TKN_AND       //36
	TKN_OR        //37
	TKN_NOT       //38
)

//homologo a la de arriba, pero para el estado
type States int

const (
	IN_START = 1 + iota
	IN_ID
	IN_NUM
	IN_LPAREN
	IN_RPAREN
	IN_LBRACE
	IN_RBRACE
	IN_SEMICOLON
	IN_ALERT
	IN_COMMA
	IN_ASSIGN
	IN_ADD
	IN_MINUS
	IN_MULTI
	IN_DIV
	IN_POW
	IN_LESS
	IN_ELESS
	IN_MORE
	IN_EMORE
	IN_EQUAL
	IN_NEQUAL
	IN_EOF
	IN_ERROR
	IN_DONE
	IN_LCOMMENT
	IN_BCOMMENT
)

//estructura donde guardaremos los tokens
type Token struct {
	tipo    token_types
	lexema  string
	linea   int
	columna int
}

var reservedWords [15]string = [15]string{
	"program", "if", "else", "fi", "do", "until", "while", "read", "write",
	"float", "int", "bool", "not", "and", "or",
}

// esta funcion lee un archivo, y regresa su contenido en un vector
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// escribe un archivo
func writeLines(lista *list.List, lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)

	for e := lista.Front(); e != nil; e = e.Next() {
		fmt.Fprintln(w, e.Value)
	}

	/*for _, line := range lines {
		fmt.Fprintln(w, line)
	}*/

	return w.Flush()
}

func isAlpha(char string) bool {
	if char >= "A" && char <= "Z" {
		return true
	} else if char >= "a" && char <= "z" {
		return true
	} else {
		return false
	}
}

func isDigit(char string, num bool, char2 string) bool {
	if char >= "0" && char <= "9" {
		return true
	} else {
		if (char == ".") && num {
			return isDigit(char2, num, char)
			//return true
		} else {
			return false
		}
	}
}

func LookUpReservedWords(lexema string, n int, i int) Token {
	var tok Token
	tok.linea = n
	tok.columna = i

	switch lexema {
	case "if":
		tok.tipo = TKN_IF
	case "read":
		tok.tipo = TKN_READ
	case "write":
		tok.tipo = TKN_WRITE
	case "fi":
		tok.tipo = TKN_FI
	case "else":
		tok.tipo = TKN_ELSE
	case "do":
		tok.tipo = TKN_DO
	case "until":
		tok.tipo = TKN_UNTIL
	case "=":
		tok.tipo = TKN_ASSIGN
	case "program":
		tok.tipo = TKN_PROGRAM
	case "float":
		tok.tipo = TKN_TIPO
	case "int":
		tok.tipo = TKN_TIPO
	case "bool":
		tok.tipo = TKN_TIPO
	case "while":
		tok.tipo = TKN_WHILE
	case "or":
		tok.tipo = TKN_OR
	case "and":
		tok.tipo = TKN_AND
	case "not":
		tok.tipo = TKN_NOT

	default:
		tok.tipo = TKN_RESERVED
	}

	for i := 0; i < 15; i++ {
		if lexema == reservedWords[i] {
			tok.lexema = lexema
			return (tok)

		}
	}

	tok.lexema = lexema
	tok.tipo = TKN_ID

	return (tok)
}

func isDelim(char string, i *int) bool {
	if char == " " {
		*i++
		return true
	} else if char == "	" {
		*i++
		return true
	} else {
		return false
	}
}

func getToken(line string, lista *list.List, blocComments *bool, n int) {
	token.lexema = ""
	token.columna = 1
	token.linea = 1
	char := ""
	tam := utf8.RuneCountInString(line)
	var i int
	casiFin := false
	var state States

	if *blocComments {
		state = IN_BCOMMENT
	} else {
		state = IN_START
	}

	for i = 0; i < tam; i++ {
		if i != tam {
			char = string([]rune(line)[i])
		}
		//fmt.Println(char)

		switch state {
		case IN_START:

			for x := i; x < tam; x++ {
				if !isDelim(char, &i) {
					x += tam
					break
				} else {
					if i < tam {
						char = string([]rune(line)[i])
					} else {
						x += tam
					}
				}

			}

			token.linea = n
			token.columna = i

			if isAlpha(char) {
				state = IN_ID
				//fmt.Println("aqui imprimo mamadas", n, i)
				token.tipo = TKN_ID
				token.lexema += char

			} else if isDigit(char, false, char) {
				state = IN_NUM
				token.tipo = TKN_NUM
				token.lexema += char

			} else if char == "(" {
				token.tipo = TKN_LPAREN
				state = IN_DONE
				token.lexema = char

			} else if char == ")" {
				token.tipo = TKN_RPAREN
				state = IN_DONE
				token.lexema = char

			} else if char == "{" {
				token.tipo = TKN_LBRACE
				state = IN_DONE
				token.lexema = char

			} else if char == "}" {
				token.tipo = TKN_RBRACE
				state = IN_DONE
				token.lexema = char

			} else if char == ";" {
				token.tipo = TKN_SEMICOLON
				state = IN_DONE
				token.lexema = char

			} else if char == "," {
				token.tipo = TKN_COMMA
				state = IN_DONE
				token.lexema = char

			} else if char == ":" {
				state = IN_ASSIGN
				token.lexema = char

			} else if char == "<" {
				state = IN_LESS
				token.lexema = char

			} else if char == ">" {
				state = IN_MORE
				token.lexema = char

			} else if char == "=" {
				state = IN_ASSIGN
				token.lexema = char

			} else if char == "+" {
				token.tipo = TKN_ADD
				state = IN_DONE
				token.lexema = char

			} else if char == "-" {
				token.tipo = TKN_MINUS
				state = IN_DONE
				token.lexema = char

			} else if char == "*" {
				token.tipo = TKN_MULTI
				state = IN_DONE
				token.lexema = char

			} else if char == "/" {
				token.tipo = TKN_MULTI
				state = IN_DIV
				token.lexema = char

			} else if char == "^" {
				token.tipo = TKN_POW
				state = IN_DONE
				token.lexema = char

			} else if char == "!" {
				token.lexema = char
				state = IN_ALERT

			} else {
				if isDelim(char, &i) {

				} else {
					fmt.Println("Error en linea ", n, " caracer no reconocido: ", char)
					token.tipo = TKN_ERROR
					token.lexema = char
					state = IN_DONE
					isCorrect = false
				}
			}

		case IN_NUM:
			char2 := "a"
			if tam > i+1 {
				char2 = string([]rune(line)[i+1])
			}
			if !isDigit(char, true, char2) {
				token.tipo = TKN_NUM
				state = IN_DONE
				i--
			} else {
				token.lexema += char
			}

		case IN_ID:
			if !((isAlpha(char) || isDigit(char, false, char)) || (char == "_")) {
				//fmt.Println("es alfanumerico Bv")
				token.tipo = TKN_ID
				state = IN_DONE
				i--
				//token.lexema[index]='\0';
				token = LookUpReservedWords(token.lexema, token.linea, token.columna)
			} else {
				token.lexema += char
			}

		case IN_ASSIGN:
			//token.lexema += char
			if char == "=" {
				token.tipo = TKN_EQUAL
				state = IN_DONE
				token.lexema += char
			} else {
				i--
				token.tipo = TKN_ASSIGN
				state = IN_DONE
			}

		case IN_LESS:
			if char == "=" {
				token.tipo = TKN_ELESS
				token.lexema += char
				state = IN_DONE
			} else {
				i--
				token.tipo = TKN_LESS
				state = IN_DONE
			}

		case IN_MORE:
			if char == "=" {
				token.tipo = TKN_EMORE
				token.lexema += char
				state = IN_DONE
			} else {
				i--
				token.tipo = TKN_MORE
				state = IN_DONE
			}

		case IN_BCOMMENT:
			if char == "*" {
				casiFin = true
			} else if char == "/" {
				if casiFin {
					state = IN_START
					*blocComments = false
					//fmt.Println("fin de comentarios de bloque")
				} else {

				}
			} else {
				if casiFin {
					casiFin = false
				}
			}

		case IN_DIV:
			if char == "*" {
				//fmt.Println("inicio de comentarios de bloque")
				state = IN_BCOMMENT
				*blocComments = true
			} else if char == "/" {
				//fmt.Println("comentarios de linea")
				i = tam
				state = IN_LCOMMENT
			} else {
				state = IN_DONE
				i--
			}

		case IN_ALERT:
			if char == "=" {
				token.tipo = TKN_NEQUAL
				token.lexema += char
				state = IN_DONE

			} else {
				token.tipo = TKN_ERROR
				state = IN_DONE
				i--
			}

		case IN_DONE:
			if token.lexema == " " || token.lexema == "" {

			} else {
				lista.PushBack(token)
				token.lexema = ""
				i--
			}
			state = IN_START
			//token.lexema = ""

		default:
			token.tipo = TKN_ERROR
			state = IN_DONE
			token.linea = n
			token.columna = i
			token.lexema += char

		}
	}
	if state == IN_LCOMMENT {

	} else if state == IN_BCOMMENT {

	} else if token.tipo == TKN_ID {
		token = LookUpReservedWords(token.lexema, token.linea, token.columna)
		lista.PushBack(token)
	} else {
		if !isDelim(token.lexema, &i) {
			if token.tipo != 0 {
				if token.lexema != "" || token.lexema != " " {
					lista.PushBack(token)
				}
			}
		}
	}

}

func Lexico(archivo string, lista *list.List) {
	/*funcion para recuperar argumentos introducidos por consola
	el primero recoje toda la linea, el segundo solo un elemento especifico
	argsWithProg := os.Args
	argsWithoutProg := os.Args[1]
	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(reservedWords)
	lines, err := readLines(argsWithoutProg)*/
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	//lista := list.New()
	blocComments := false

	for i, line := range lines {
		getToken(line, lista, &blocComments, i+1)
		//fmt.Println(i, line)
	}

	err = writeLines(lista, lines, "tokens.txt")
}
