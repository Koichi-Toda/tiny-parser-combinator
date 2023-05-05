package main

import (
	"github.com/Koichi-Toda/tiny-parser-combinator/parser"
)

func main() {
	expr := parser.Parser_sample()
	parser.ParseAll(expr, "(1+2+3)*(4+5+6)")
}
