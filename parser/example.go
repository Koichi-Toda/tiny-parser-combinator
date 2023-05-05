package parser

import (
	"regexp"
	"strconv"
)

func Parser_sample() Parser {
	// number regexp
	numrex := regexp.MustCompile("^[1-9][0-9]*")

	// basic matching function
	p := func(kind string, in string) (bool, Token, string) {
		ilen, klen := len(in), len(kind)
		if ilen < klen {
			return false, Token{in}, in
		}
		if item := in[:klen]; item == kind {
			return true, Token{kind}, in[klen:]
		} else {
			return false, Token{item}, in
		}
	}

	// number
	num := func(kind string, in string) (bool, Token, string) {
		ret := numrex.FindStringIndex(in)
		if ret == nil {
			eletter := in[:1]
			return false, Token{eletter}, in
		}
		index := ret[1]

		number, _ := strconv.Atoi(in[:index])
		return true, Token{number}, in[index:]
	}

	// generate operation fuction
	f1 := func(t []Token) []Token {
		var opfunc func(int) int
		y := t[1].Value.(int)

		switch t[0].String() {
		case "+":
			opfunc = func(x int) int { return x + y }
		case "-":
			opfunc = func(x int) int { return x - y }
		case "*":
			opfunc = func(x int) int { return x * y }
		case "/":
			opfunc = func(x int) int { return x / y }
		}
		return []Token{{opfunc}}
	}

	// fold functions
	f2 := func(t []Token) []Token {
		sum := t[0].Value.(int)
		for _, item := range t[1:] {
			sum = item.Value.(func(int) int)(sum)
		}
		return []Token{{sum}}
	}

	// reject paren caracter
	f3 := func(t []Token) []Token {
		// reject "(" and ")"
		// t[0]("("),t[1](value),t[2](")") -> t[1](value)
		ret := []Token{t[1]}
		return ret
	}

	var expr, term, factor Parser
	e := func(kind string) Parser { return Elem(kind, p) }
	number := Elem("Number", num)

	// parser rule (basic calculation rule: Four arithmetic operations)
	expr = T_(Seq(&term,
		Rep(Choice(
			T_(Seq(e("+"), &term), f1),
			T_(Seq(e("-"), &term), f1)))), f2)

	term = T_(Seq(&factor,
		Rep(Choice(
			T_(Seq(e("*"), &factor), f1),
			T_(Seq(e("/"), &factor), f1)))), f2)
	factor = Choice(
		number,
		T_(Seq(e("("), Seq(&expr, e(")"))), f3))

	return expr
}
