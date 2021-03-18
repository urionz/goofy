package interpreter

import (
	"fmt"
	"go/parser"
	"go/scanner"
	"go/token"
)

type Type uint8

const (
	Import = Type(iota)
	TypeDecl
	FuncDecl
	VarDecl
	Shell
	Print
	Unknown
	Expr
	Stmt
	Empty
)

func createScannerFor(code string) scanner.Scanner {
	var s scanner.Scanner
	fs := token.NewFileSet()
	s.Init(fs.AddFile("", fs.Base(), len(code)), []byte(code), nil, scanner.ScanComments)
	return s
}
func tokenizerAndLiterizer(code string) ([]token.Token, []string) {
	s := createScannerFor(code)
	var tokens []token.Token
	var lits []string
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		tokens = append(tokens, tok)
		if tok == token.IDENT || tok == token.STRING {
			lits = append(lits, lit)
		} else if tok == token.INT || tok == token.FLOAT {
			lits = append(lits, fmt.Sprint(code[pos-1:pos]))
		} else {
			lits = append(lits, tok.String())
		}
	}
	return tokens, lits
}

func Parse(code string) (Type, error) {
	if isEmpty(code) {
		return Empty, nil
	} else if isShellCommand(code) {
		return Shell, nil
	} else if isImport(code) {
		return Import, nil
	} else if IsFuncDecl(code) {
		return FuncDecl, nil
	} else if isTypeDecl(code) {
		return TypeDecl, nil
	} else if isPrint(code) {
		return Print, nil
	} else if isVarDecl(code) {
		return VarDecl, nil
	} else if isExpr(code) {
		return Expr, nil
	} else {
		return Unknown, nil
	}
}

func isExpr(code string) bool {
	_, err := parser.ParseExpr(code)
	if err != nil {
		return false
	}
	return true
}

func ShouldContinue(code string) (int, bool) {
	var stillOpenChars int
	for _, c := range code {
		if c == '{' || c == '(' {
			stillOpenChars++
			continue
		}

		if c == '}' || c == ')' {
			stillOpenChars--
		}
	}
	return stillOpenChars, stillOpenChars > 0
}
func isEmpty(code string) bool {
	return len(code) == 0
}

func isPrint(code string) bool {
	tokens, lits := tokenizerAndLiterizer(code)
	for i, t := range tokens {
		if t == token.IDENT &&
			(lits[i] == "Println" || lits[i] == "Printf" || lits[i] == "Print" || lits[i] == "println") || lits[i] == "printf" || lits[i] == "print" {
			return true
		}
	}
	return false
}
