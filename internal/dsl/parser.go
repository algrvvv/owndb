package dsl

type Parser struct {
	tokens    []Token
	pos       int
	DebugMode bool
}

func NewParser(tokens []Token) Parser {
	return Parser{
		tokens: tokens,
	}
}
