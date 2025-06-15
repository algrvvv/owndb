package dsl

type TokenType int

type Token struct {
	Type     TokenType
	Position Position
	Value    string
}

// nolint
//
// NOTE: отключаю линтер, так как жалуется на андескор и капс в именах.
// а я ведь просто вдохнавился go/token
const (
	// special
	EOF TokenType = iota
	ILLEGAL

	// literals
	lit_start
	IDENT   // "main"
	KEYWORD // keyword mark
	TYPE    // type mark
	SEMI    // ;
	lit_end

	types_start
	INT    // 123
	FLOAT  // 123.123
	STRING // "string"
	BOOL   // bool
	types_end

	// operators
	operators_start
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ // !=
	LEQ // <=
	GEQ // >=
	operators_end

	// keywords
	keywords_start
	SET    // set key = value
	GET    // get key
	RM     // rm key
	IF     // if key = value then
	THEN   // if key = value then
	EXISTS // if EXISTS key then
	SAVE   // save
	KEYS   // show all availble keys
	SHOW   // show all data
	EXP    // explain query (debug mode)
	keywords_end

	booltype_start
	TRUE
	FALSE
	booltype_end
)

var tokens = [...]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	KEYWORD: "KEYWORD",
	TYPE:    "TYPE",
	INT:     "INT",
	FLOAT:   "FLOAT",
	BOOL:    "BOOL",
	STRING:  "STRING",
	SEMI:    ";",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ: "!=",
	LEQ: "<=",
	GEQ: ">=",

	SET:    "SET",
	GET:    "GET",
	RM:     "RM",
	IF:     "IF",
	THEN:   "THEN",
	EXISTS: "EXISTS",
	SAVE:   "SAVE",
	KEYS:   "KEYS",
	SHOW:   "SHOW",
	EXP:    "EXP",

	TRUE:  "TRUE",
	FALSE: "FALSE",
}

var (
	keywords map[string]TokenType
	literals map[string]TokenType
	types    map[string]TokenType
)

func init() {
	keywords = make(map[string]TokenType, keywords_end-(keywords_start+1))
	for i := keywords_start + 1; i < keywords_end; i++ {
		keywords[tokens[i]] = i
	}

	literals = make(map[string]TokenType, lit_end-(lit_start+1))
	for i := lit_start + 1; i < lit_end; i++ {
		literals[tokens[i]] = i
	}

	types = make(map[string]TokenType, types_end-(types_start+1))
	for i := types_start + 1; i < types_end; i++ {
		types[tokens[i]] = i
	}
}

func (t TokenType) IsType() bool { return t > types_start && t < types_end }

func (t TokenType) IsLiteral() bool { return t > lit_start && t < lit_end }

func (t TokenType) String() string { return tokens[t] }

func (t TokenType) IsKeyword() bool { return t > keywords_start && t < keywords_end }

func (t TokenType) IsOperator() bool { return t > operators_start && t < operators_end }

func IsBool(lit string) bool {
	return lit == tokens[TRUE] || lit == tokens[FALSE]
}

func LookupTokTyp(lit string) TokenType {
	for k, v := range tokens {
		if v == lit {
			return TokenType(k)
		}
	}

	return -1
}

func Lookup(lit string) TokenType {
	if _, isKeyword := keywords[lit]; isKeyword {
		return KEYWORD
	}

	if lit, isLit := literals[lit]; isLit {
		return lit
	}

	if _, isType := types[lit]; isType {
		return TYPE
	}

	return IDENT
}
