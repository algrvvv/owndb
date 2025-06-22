package exec

import (
	"fmt"
	"strings"

	"github.com/algrvvv/owndb/internal/dsl"
)

type InlineDiscardExecutor struct {
	intr *dsl.Interpreter
}

func NewInlineDiscardExecutor(intr *dsl.Interpreter) Executor {
	return &InlineDiscardExecutor{
		intr: intr,
	}
}

func (i *InlineDiscardExecutor) Execute(command string) (res any, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
	}()

	var tokens []dsl.Token
	lexer := dsl.NewLexer(strings.NewReader(command))

	for {
		pos, tok, lit := lexer.Lex()
		if tok == dsl.EOF {
			break
		}

		tokens = append(tokens, dsl.Token{Type: tok, Position: pos, Value: lit})
	}

	parser := dsl.NewParser(tokens)
	stmt, err := parser.Parse()
	if err != nil || stmt == nil {
		return
	}

	if parser.DebugMode {
		for _, tok := range tokens {
			fmt.Printf("%d:%d\t%s\t%s\n", tok.Position.Line, tok.Position.Column, tok.Type, tok.Value)
		}
	}

	res, err = i.intr.ExecStatement(stmt, parser.DebugMode)
	return
}
