package dsl_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/algrvvv/owndb/internal/dsl"
)

const testsHashFile = "lexer_test_hash.json"

func TestLexer(t *testing.T) {
	var rewriteHash bool
	args := os.Args[1:]

	for _, arg := range args {
		if arg == "rewrite" {
			rewriteHash = true
		}
	}

	hashes, err := readHash()
	if err != nil {
		panic(err)
	}

	tests := []struct {
		command string
	}{
		{
			command: "SET key = value;",
		},
		{
			command: "SET INT another = 123321;",
		},
		{
			command: "GET key;",
		},
		{
			command: "SET STRING foo = bar;",
		},
		{
			command: "SET BOOL bar = false;",
		},
		{
			command: "RM bar;",
		},
		{
			command: "SAVE;",
		},
		{
			command: "IF another = 123321 THEN GET bar;",
		},
		{
			command: "IF key = foo THEN SET INT bar = 12;",
		},
		{
			command: "IF key = foo THEN SET INT bar = 12;",
		},
		{
			command: "SET FLOAT val = 1.23;",
		},
	}

	m := make(map[string]string, len(tests))
	var errorCounter int
	for _, tt := range tests {
		println()
		fmt.Println("input: ", tt.command)
		reader := strings.NewReader(tt.command)
		lexer := dsl.NewLexer(reader)

		var sum string
		for {
			pos, tok, lit := lexer.Lex()
			if tok == dsl.EOF {
				break
			}

			line := fmt.Sprintf("%d:%d\t%s\t%s", pos.Line, pos.Column, tok, lit)
			sum += line + "\n"

			fmt.Println(line)
		}

		h := hash(sum)
		wantedHash := hashes[tt.command]

		if h != wantedHash {
			t.Errorf("FAIL! failed to lex: %q\nlex result:\n%s\nwant hash: %s\ngot hash: %s", tt.command, sum, wantedHash, h)
			errorCounter++
		} else {
			println("\nstatus: OK")
		}

		m[tt.command] = h
		fmt.Println("hash: ", h)
		println()
	}

	if rewriteHash {
		if err := writeHash(m); err != nil {
			panic(err)
		} else {
			println("OK! tests hash rewrited")
		}
	}

	if errorCounter == 0 {
		fmt.Println("OK! tests passed")
	} else {
		fmt.Println("FAIL! errors: ", errorCounter)
	}
}

func writeHash(data map[string]string) error {
	jsonData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(testsHashFile, jsonData, 0600)
}

func readHash() (map[string]string, error) {
	data, err := os.ReadFile(testsHashFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, nil
		}
	}

	var m map[string]string
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func hash(str string) string {
	h := sha256.Sum256([]byte(str))
	return hex.EncodeToString(h[:])
}
