package dsl

import "fmt"

// peek метод, который возвращает токен по позиции парсера без передвижения
func (p *Parser) peek() Token {
	return p.tokens[p.pos]
}

func (p *Parser) next() Token {
	if p.pos >= len(p.tokens) {
		panic("token out of bounds")
	}

	p.pos++
	return p.tokens[p.pos]
}

func (p *Parser) nextExpect(typ TokenType) Token {
	p.pos++

	if p.pos >= len(p.tokens) {
		prevtok := p.tokens[p.pos-1]
		panic(fmt.Sprintf("expect %q token on: %d:%d;", typ.String(), prevtok.Position.Line, prevtok.Position.Column+1))
	}

	tok := p.tokens[p.pos]
	if tok.Type != typ {
		panic(fmt.Sprintf("next token have unexpected type. got %q instead of %q;", tok.Type.String(), typ.String()))
	}

	return tok
}

// advance метод, который возращает токен по позиции парсера с передвижением
func (p *Parser) advance() Token {
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

// match метод, который проверяем тип токена и сам токен
func (p *Parser) match(typ TokenType, valTyp TokenType) bool {
	if p.pos >= len(p.tokens) {
		return false
	}

	tok := p.peek()
	val := tokens[valTyp]
	return typ == tok.Type && val == tok.Value
}

func (p *Parser) expect(typ TokenType) bool {
	if p.pos >= len(p.tokens) {
		return false
	}

	tok := p.peek()
	return typ == tok.Type
}
