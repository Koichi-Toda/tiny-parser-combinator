package parser

import (
	"fmt"
	"strconv"
)

/*
tiny-parser-combinator ... parser combinator sample code

This code is based on "Programming in Scala. 4th Edition.",
Especially Chapter 33[Combinator Parsing] section.


parser function list
(scala version)	| (this program)
p~q				: Seq(p,q)
p<~q			: SeqLeft(p,q)
p~>q			: SeqRight(p,q)
p|q				: Choice(p,q)
opt(p)			: Opt(p)
rep(p)			: Rep(p)
repsep(p,q)		: RepSep(p,q)
p ^^ f			: T_(p,f)
*/

type Parser func(string) ParseResult

type Token struct {
	Value any
}

func (t Token) String() string {
	switch value := t.Value.(type) {
	case string:
		return value
	case int:
		return strconv.Itoa(value)
	}

	return ""
}

type ParseResult interface {
	String() string
}

type Success struct {
	result string
	token  []Token
	in     string
}

type Failure struct {
	result string
	token  []Token
	in     string
}

func (s Success) String() string {
	return fmt.Sprintf("[Parse Success] accepted:%v,result:%v, rest:%v",
		s.result, s.token, s.in)
}

func (f Failure) String() string {
	return fmt.Sprintf("[Parse Failure] message:%v, result:%v, rest:%v",
		f.result, f.token, f.in)
}

// p_change(parser any) returns Parser.
// The reason why this function is necessary is that golang cannot represent function pointers.
// (We cannot use &seq(&rep(&choice(...))) in golang.)
// On the other hand, in order to make circular references to expr, term, and functor,
// the function must be passed by pointer.
// if every parser function use pointer, we need to use each function assigned to each variable.
//
//	c = choice(...); r = rep(&c); seq(&r)
//
// but I want to express parser combinator as follows.
//
//	seq(rep(choice(...)))
//
// In this time, I decided to use any(=interface{}).
// This allows us to pass both pointers and non-pointers.
// The trade-off is that we have to do some type checking and mending later,
// such as the following.
func p_change(parser any) Parser {
	switch p := parser.(type) {
	case Parser:
		return p
	case *Parser:
		return (*p)
	default:
		panic(fmt.Errorf("error type '%T' ", parser))
	}
}

// elemental parsing
func Elem(kind string, p func(string, string) (bool, Token, string)) Parser {
	return func(in string) ParseResult {
		if r, token, rest := p(kind, in); r {
			return Success{token.String(), []Token{token}, rest}
		} else {
			return Failure{fmt.Sprintf("%v expected but %v found.", kind, token), []Token{token}, in}
		}
	}
}

// sequential parsing
func Seq(p_, q_ any) Parser {
	return func(in string) ParseResult {
		p, q := p_change(p_), p_change(q_)

		switch r := p(in).(type) {
		case Success:
			switch r2 := q(r.in).(type) {
			case Success:
				return Success{
					r.result + r2.result,
					append(r.token, r2.token...),
					r2.in,
				}
			case Failure:
				return r2
			}
		case Failure:
			return r
		}

		return Failure{"Unexpeced error happen!", []Token{}, in}
	}
}

// alternative parsing
func Choice(p Parser, q Parser) Parser {
	return func(in string) ParseResult {
		r := p(in)
		switch r.(type) {
		case Success:
			return r
		case Failure:
			return q(in)
		}

		return Failure{"Unexpeced error happen!", []Token{}, in}
	}
}

// result transformation (token -> f(token))
func T_(p Parser, f func([]Token) []Token) Parser {
	return func(in string) ParseResult {
		r := p(in)
		switch v := r.(type) {
		case Success:
			return Success{v.result, f(v.token), v.in}
		case Failure:
			return r
		}

		return Failure{"Unexpeced error happen!", []Token{}, in}
	}
}

func SeqLeft(p, q Parser) Parser {
	return T_(Seq(p, q), func(t []Token) []Token {
		return t[:1]
	})
}

func SeqRight(p, q Parser) Parser {
	return T_(Seq(p, q), func(t []Token) []Token {
		return t[1:]
	})
}

func success() Parser {
	return func(in string) ParseResult {
		return Success{"", nil, in}
	}
}

func failure() Parser {
	return func(in string) ParseResult {
		return Failure{"", nil, in}
	}
}

// optional parsing
func Opt(p Parser) Parser {
	return func(in string) ParseResult {
		opt := Choice(p, success())
		return opt(in)
	}
}

// repeatable parsing
func Rep(p Parser) Parser {
	return func(in string) ParseResult {
		rep := Choice(Seq(p, Rep(p)), success())
		return rep(in)
	}
}

// repeatable parsing with separator
func RepSep(p, q Parser) Parser {
	return func(in string) ParseResult {
		repsep := Seq(p, Rep(SeqRight(q, p)))
		return repsep(in)
	}
}

func ParseAll(p Parser, input string) {
	result := p(input)
	fmt.Printf("parsed: %v\n", result)
}
