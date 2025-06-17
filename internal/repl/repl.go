package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/algrvvv/owndb/internal/dsl"
)

type REPL struct {
	intr *dsl.Interpreter
}

func NewREPLInstance(intr *dsl.Interpreter) *REPL {
	return &REPL{
		intr: intr,
	}
}

func (r *REPL) Scan(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			print("> ")
			continue
		} else if line == "EXIT" || line == "EXIT;" {
			return fmt.Errorf("exit by user")
		}

		data, err := r.ExecQuery(line)
		if err != nil {
			fmt.Println("FAIL! ", err.Error())
			print("\n> ")
			continue
		}

		if data != nil {
			fmt.Println(data)
		}

		print("\n> ")
	}

	if scanner.Err() != nil {
		return fmt.Errorf("scanner error")
	}

	return nil
}

func (r *REPL) ExecQuery(query string) (res any, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
	}()

	var tokens []dsl.Token
	lexer := dsl.NewLexer(strings.NewReader(query))

	for {
		pos, tok, lit := lexer.Lex()
		if tok == dsl.EOF {
			break
		}

		tokens = append(tokens, dsl.Token{Type: tok, Position: pos, Value: lit})
	}

	parser := dsl.NewParser(tokens)
	stmt, err := parser.Parse()
	if err != nil {
		return
	}

	if parser.DebugMode {
		for _, tok := range tokens {
			fmt.Printf("%d:%d\t%s\t%s\n", tok.Position.Line, tok.Position.Column, tok.Type, tok.Value)
		}
	}

	res, err = r.intr.ExecStatement(stmt, parser.DebugMode)
	return
}
