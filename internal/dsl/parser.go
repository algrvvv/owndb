package dsl

import (
	"fmt"
	"strconv"
)

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

func (p *Parser) Parse() (Statement, error) {
	// TODO:
	// + 1. берем следующий токен.
	// + 2. проверяем, чтобы выражение содержало в себе как минимум
	// одно ключевое слово: к примеру, SET, GET и тд.
	// + 3. запускаем парсер для найденной команды

	switch {
	case p.match(KEYWORD, SET):
		return p.parseSet()
	case p.match(KEYWORD, GET):
		return p.parseGet()
	case p.match(KEYWORD, SHOW):
		return p.parseShow()
	case p.match(KEYWORD, KEYS):
		return p.parseKeys()
	case p.match(KEYWORD, RM):
		return p.parseRemove()
	case p.match(KEYWORD, SAVE):
		return p.parseSave()
	case p.match(KEYWORD, EXP):
		p.pos++
		p.DebugMode = true
		return p.Parse()
	}

	return nil, fmt.Errorf("invalid keyword statement")
}

func (p *Parser) parseExpression(tok Token) (Expression, error) {
	switch tok.Type {
	case IDENT:
		return &Identifier{Name: tok.Value}, nil
	case INT:
		v, err := strconv.Atoi(tok.Value)
		if err != nil {
			return nil, err
		}

		return &Literal[int]{Value: v}, nil
	case FLOAT:
		v, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, err
		}

		return &Literal[float64]{Value: v}, nil
	case STRING:
		return &Literal[string]{Value: tok.Value}, nil
	case BOOL:
		v, err := strconv.ParseBool(tok.Value)
		if err != nil {
			return nil, err
		}

		return &Literal[bool]{Value: v}, nil
	default:
		panic("unexpected expresssion token type")
	}
}

// parseSet метод для парсинга команды на добавление данных.
//
// пример запроса на добавление данных:
// SET ?TYPE key = val;
//
// в случае, запроса на сохранение данных у нас может быть,
// а может и не быть добавлен тип данных.
// по умолчанию его ставим, как STRING
func (p *Parser) parseSet() (Statement, error) {
	// NOTE: пробуем дернуть тип данных,
	// если это не тип данных, то тип данных ставим, как стринг,
	// а полученный токен используем, как ключ

	var typ TokenType
	var name Token

	tok := p.next()

	// пытаемся получить тип переменной
	if p.expect(TYPE) {
		// если токен является типом данных, то сохраняем его,
		// а также читаем следующий токен, который должен быть
		// названием ключа
		typ = LookupTokTyp(tok.Value)
		name = p.nextExpect(IDENT)
	} else {
		// если мы видим, что это не тип данных, то в таком
		// случае, это должно быть название ключа,
		// поэтому проверяем тип токена и если все правильно,
		// то сохраняем название ключа.
		//
		// и если это оказался не тип данных,
		// то выбираем тип данных по умолчанию,
		// то есть стринга

		typ = STRING
		if p.expect(IDENT) {
			name = tok
		}
	}

	p.nextExpect(ASSIGN)                 // пропускаем ИМЕННО равно
	v := p.next()                        // получаем значение. не используем nextExpect, так как там могут быть разные типы данных
	valExpr, err := p.parseExpression(v) // получаем выражение для значения
	if err != nil {
		return nil, err
	}

	switch valExpr.(type) {
	case *Literal[string]:
		if typ != STRING {
			panic(fmt.Sprintf("unexpected value type. expected: %q, got: %q", typ.String(), STRING.String()))
		}
	case *Literal[int]:
		if typ != INT {
			panic(fmt.Sprintf("unexpected value type. expected: %q, got: %q", typ.String(), INT.String()))
		}
	case *Literal[float64]:
		if typ != FLOAT {
			panic(fmt.Sprintf("unexpected value type. expected: %q, got: %q", typ.String(), FLOAT.String()))
		}
	case *Literal[bool]:
		if typ != BOOL {
			panic(fmt.Sprintf("unexpected value type. expected: %q, got: %q", typ.String(), BOOL.String()))
		}
	}

	p.nextExpect(SEMI) // проверяем, что наша команда закрывается

	return &SetStatement{
		Type:  typ,
		Name:  name.Value,
		Value: valExpr,
	}, nil
}

// parseGet метод, который создает инструкцию для команды
// получения данных.
//
// обычно запрос будет выглядеть так
// GET key;
//
// поэтому, в этом случае, достаточно проверить, что
// следующий токен является ключом, а так же,
// что запрос закрывается.
func (p *Parser) parseGet() (Statement, error) {
	tok := p.nextExpect(IDENT)
	p.nextExpect(SEMI)

	return &GetStatement{
		Name: tok.Value,
	}, nil
}

func (p *Parser) parseShow() (stmt Statement, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", err)
		}
	}()

	limit := -1
	tok := p.next()

	if tok.Type != INT && tok.Type != SEMI {
		return nil, fmt.Errorf("invalid limit param type: %s", tok.Type.String())
	} else if tok.Type == INT {
		limit, err = strconv.Atoi(tok.Value)
		if err != nil {
			return
		}
	}

	return &ShowStatement{
		Limit: limit,
	}, nil
}

func (p *Parser) parseKeys() (stmt Statement, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", err)
		}
	}()

	limit := -1
	tok := p.next()

	if tok.Type != INT && tok.Type != SEMI {
		return nil, fmt.Errorf("invalid limit parametr type: %s", tok.Type.String())
	} else if tok.Type == INT {
		limit, err = strconv.Atoi(tok.Value)
		if err != nil {
			return
		}
	}

	return &KeysStatement{
		Limit: limit,
	}, nil
}

func (p *Parser) parseRemove() (Statement, error) {
	tok := p.next()
	p.nextExpect(SEMI)

	return &RemoveStatement{
		Name: tok.Value,
	}, nil
}

func (p *Parser) parseSave() (Statement, error) {
	p.nextExpect(SEMI)
	return &SaveStatement{}, nil
}
