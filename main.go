package main

import (
	"fmt"
	"strings"
)

// Token merepresentasikan hasil dari analisis leksikal
type Token struct {
	Type  string
	Value string
}

// FuncDeclCompiler adalah struct utama untuk proses kompilasi deklarasi fungsi
type FuncDeclCompiler struct {
	sourceCode  string
	tokens      []Token
	pos         int
	tempCounter int
}

func NewFuncDeclCompiler(source string) *FuncDeclCompiler {
	return &FuncDeclCompiler{
		sourceCode:  source,
		tempCounter: 1,
	}
}

// 1. Analisis Leksikal
func (c *FuncDeclCompiler) LexicalAnalysis() []Token {
	// Memisahkan simbol-simbol khusus agar mudah dipecah (tokenizing)
	replacements := map[string]string{
		"(": " ( ",
		")": " ) ",
		"{": " { ",
		"}": " } ",
		",": " , ",
		"=": " = ",
		"+": " + ",
		"-": " - ",
		"*": " * ",
		"/": " / ",
	}

	src := c.sourceCode
	for k, v := range replacements {
		src = strings.ReplaceAll(src, k, v)
	}

	// Memecah berdasarkan spasi
	words := strings.Fields(src)
	var tokens []Token

	for _, word := range words {
		tokens = append(tokens, Token{Type: "TOKEN", Value: word})
	}
	c.tokens = tokens
	return tokens
}

// ASTNode merepresentasikan simpul pada Abstract Syntax Tree
type ASTNode struct {
	Type     string
	Value    string
	Children []*ASTNode
}

// 2. Analisis Sintaksis
func (c *FuncDeclCompiler) SyntaxAnalysis() *ASTNode {
	c.pos = 0
	if c.match("func") {
		name := c.consume()

		c.consumeExpected("(")
		var params []string
		for c.peek() != ")" {
			paramName := c.consume()
			paramType := c.consume()
			params = append(params, paramName+" "+paramType)
			if c.peek() == "," {
				c.consumeExpected(",")
			}
		}
		c.consumeExpected(")")

		returnType := c.consume()

		c.consumeExpected("{")
		var body []string
		for c.peek() != "}" {
			body = append(body, c.consume())
		}
		c.consumeExpected("}")

		// Membuat struktur AST
		funcNode := &ASTNode{Type: "FuncDecl", Value: name}
		paramsNode := &ASTNode{Type: "Params", Value: strings.Join(params, ", ")}
		returnTypeNode := &ASTNode{Type: "ReturnType", Value: returnType}
		bodyNode := &ASTNode{Type: "Body", Value: strings.Join(body, " ")}

		funcNode.Children = append(funcNode.Children, paramsNode, returnTypeNode, bodyNode)
		return funcNode
	}
	panic("Syntax Error: Diharapkan kata kunci 'func'")
}

// Helper methods untuk proses parsing
func (c *FuncDeclCompiler) peek() string {
	if c.pos < len(c.tokens) {
		return c.tokens[c.pos].Value
	}
	return ""
}

func (c *FuncDeclCompiler) match(val string) bool {
	if c.peek() == val {
		c.pos++
		return true
	}
	return false
}

func (c *FuncDeclCompiler) consume() string {
	val := c.peek()
	c.pos++
	return val
}

func (c *FuncDeclCompiler) consumeExpected(val string) {
	if !c.match(val) {
		panic(fmt.Sprintf("Syntax Error: Diharapkan '%s', tapi mendapatkan '%s'", val, c.peek()))
	}
}

// 3. Analisis Semantik
func (c *FuncDeclCompiler) SemanticAnalysis(ast *ASTNode) {
	// Pengecekan dasar: validasi tipe kembalian dan tipe parameter
	if ast.Type != "FuncDecl" {
		panic("Semantic Error: Node AST tidak valid")
	}

	paramsNode := ast.Children[0]
	returnTypeNode := ast.Children[1]

	validTypes := map[string]bool{"int": true, "float": true, "string": true, "void": true}

	// Cek validitas return type
	if !validTypes[returnTypeNode.Value] {
		panic(fmt.Sprintf("Semantic Error: Tipe pengembalian tidak valid '%s'", returnTypeNode.Value))
	}

	// Cek validitas parameter
	if paramsNode.Value != "" {
		paramList := strings.Split(paramsNode.Value, ", ")
		for _, param := range paramList {
			parts := strings.Split(param, " ")
			if len(parts) == 2 {
				pType := parts[1]
				if !validTypes[pType] {
					panic(fmt.Sprintf("Semantic Error: Tipe parameter tidak valid '%s'", pType))
				}
			}
		}
	}
	fmt.Println("Semantic Analysis Passed: Tipe data parameter dan return type valid.")
}

// 4. Generasi Kode Antara (Three-Address Code / TAC)
func (c *FuncDeclCompiler) GenerateTAC(ast *ASTNode) string {
	var tac []string

	name := ast.Value
	params := ast.Children[0].Value
	bodyStr := ast.Children[2].Value

	tac = append(tac, fmt.Sprintf("BeginFunc %s", name))

	// Deklarasi parameter
	if params != "" {
		paramList := strings.Split(params, ", ")
		for _, param := range paramList {
			parts := strings.Split(param, " ")
			tac = append(tac, fmt.Sprintf("PopParam %s", parts[0]))
		}
	}

	// Mengurai dan mengonversi body (statement) menjadi TAC
	bodyTokens := strings.Fields(bodyStr)
	for i := 0; i < len(bodyTokens); i++ {
		if bodyTokens[i] == "return" {
			if i+1 < len(bodyTokens) {
				tac = append(tac, fmt.Sprintf("Return %s", bodyTokens[i+1]))
			} else {
				tac = append(tac, "Return")
			}
		} else if i+1 < len(bodyTokens) && bodyTokens[i+1] == "=" {
			// Sederhana, meng-handle format: a = b + c
			if i+4 < len(bodyTokens) && (bodyTokens[i+3] == "+" || bodyTokens[i+3] == "-" || bodyTokens[i+3] == "*" || bodyTokens[i+3] == "/") {
				tac = append(tac, fmt.Sprintf("t%d = %s %s %s", c.tempCounter, bodyTokens[i+2], bodyTokens[i+3], bodyTokens[i+4]))
				tac = append(tac, fmt.Sprintf("%s = t%d", bodyTokens[i], c.tempCounter))
				c.tempCounter++
				i += 4
			} else {
				// Handle: a = b
				tac = append(tac, fmt.Sprintf("%s = %s", bodyTokens[i], bodyTokens[i+2]))
				i += 2
			}
		}
	}

	tac = append(tac, fmt.Sprintf("EndFunc %s", name))

	return strings.Join(tac, "\n")
}

func main() {
	// Contoh source code: Deklarasi Fungsi / Metode
	source := "func hitungTotal ( a int , b int ) int { c = a + b return c }"
	compiler := NewFuncDeclCompiler(source)

	fmt.Println("--- 1. Analisis Leksikal (Tokens) ---")
	tokens := compiler.LexicalAnalysis()
	for _, t := range tokens {
		fmt.Printf("'%s' ", t.Value)
	}

	fmt.Println("--- 2. Analisis Sintaksis (Abstract Syntax Tree) ---")
	ast := compiler.SyntaxAnalysis()
	fmt.Printf("FuncDecl: %s\n", ast.Value)
	fmt.Printf("  ├─ Params: %s\n", ast.Children[0].Value)
	fmt.Printf("  ├─ ReturnType: %s\n", ast.Children[1].Value)
	fmt.Printf("  └─ Body: %s\n\n", ast.Children[2].Value)

	fmt.Println("--- 3. Analisis Semantik ---")
	compiler.SemanticAnalysis(ast)
	fmt.Println()

	fmt.Println("--- 4. Generasi Kode Antara (TAC) ---")
	tac := compiler.GenerateTAC(ast)
	fmt.Println(tac)
}
