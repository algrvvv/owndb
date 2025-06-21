package exec

import (
	"fmt"
	"strings"

	"github.com/algrvvv/owndb/internal/dsl"
	"github.com/algrvvv/owndb/internal/wal"
)

type InlineExecutor struct {
	intr *dsl.Interpreter
	wal  *wal.WAL
}

func NewInlineExecutor(intr *dsl.Interpreter, wal *wal.WAL) Executor {
	return &InlineExecutor{
		intr: intr,
		wal:  wal,
	}
}

func (i *InlineExecutor) Execute(command string) (res any, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
	}()

	err = i.wal.Write(command)
	if err != nil {
		return nil, err
	}

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
