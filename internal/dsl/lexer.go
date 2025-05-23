package dsl

import (
	"bufio"
	"errors"
	"io"
	"unicode"
)

type Position struct {
	Line   int
	Column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

// NewLexer returns new lexer instance.
func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{Line: 1, Column: 0},
		reader: bufio.NewReader(reader),
	}
}

// Lex scans the input for the next token.
// it returns position, token and lit.
func (l *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return l.pos, EOF, ""
			}

			// panic(err)
			return l.pos, ILLEGAL, ""
		}

		l.pos.Column++

		// TODO: дописать логику проверки !=, ==, <=, >=
		switch r {
		case '\n':
			l.resetPosition()
		case ';':
			return l.pos, SEMI, ";"
		case '+':
			return l.pos, ADD, "+"
		case '-':
			return l.pos, SUB, "-"
		case '*':
			return l.pos, MUL, "*"
		case '/':
			return l.pos, DIV, "/"
		case '=':
			return l.pos, ASSIGN, "="
		case '!':
			return l.pos, NOT, "!"
		case '<':
			return l.pos, LSS, "<"
		case '>':
			return l.pos, GTR, ">"
		default:
			if unicode.IsSpace(r) {
				// if it is space just continue
				continue
			} else if unicode.IsDigit(r) {
				startPos := l.pos
				l.backup()
				lit, tok := l.lexNum()
				return startPos, tok, lit
			} else if unicode.IsLetter(r) {
				startPos := l.pos
				l.backup()
				lit := l.lexStr()

				tok := Lookup(lit)
				return startPos, tok, lit
			}
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.Line++
	l.pos.Column = 0
}

// backup метод, которые возвращает последний прочитанный элементв буфер.
func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.Column--
}

func (l *Lexer) lexNum() (string, Token) {
	var tok Token = INT
	var lit []rune

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return string(lit), tok
			}

			panic(err)
		}

		l.pos.Column++

		if unicode.IsDigit(r) {
			// если это число, то добавляем к прошлым рунам
			lit = append(lit, r)
		} else if r == '.' || r == ',' {
			// здесь мы по факту можем получить float
			// который будет исползовать . или ,
			// в таком случае, мы просто игнорируем и проверяем следующий символ
			// если он является числом, то кладем его и идем дальше.
			// если это не число, то кладем обратно уже два символа и возвращем то, что уже запарсили

			nextr, _, err := l.reader.ReadRune()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return string(lit), tok
				}

				panic(err)
			}

			if unicode.IsDigit(nextr) {
				lit = append(lit, '.')   // кладем сначала разделитель
				lit = append(lit, nextr) // кладем следующее число, чтобы не читать его еще раз
				tok = FLOAT              // меняем токен на флот
			} else {
				l.backup() // возращаем следующий символ
				l.backup() // возращаем разделитель
			}
		} else {
			// если это уже не число, то откатываемся на один символ назад
			// и выходим из цикла
			l.backup()
			break
		}
	}

	return string(lit), tok
}

func (l *Lexer) lexStr() string {
	var lit []rune

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return string(lit)
			}

			panic(err)
		}

		l.pos.Column++
		if unicode.IsLetter(r) {
			lit = append(lit, r)
		} else {
			l.backup()
			return string(lit)
		}
	}
}
