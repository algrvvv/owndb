package dsl

type Token int

// nolint
//
// NOTE: отключаю линтер, так как жалуется на андескор и капс в именах.
// а я ведь просто вдохнавился go/token
const (
	// special
	EOF Token = iota
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
	keywords_end
)

var tokens = [...]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	KEYWORD: "KEYWORD",
	TYPE:    "TYPE",
	INT:     "INT",
	FLOAT:   "FLOAT",
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
}

var (
	keywords map[string]Token
	literals map[string]Token
	types    map[string]Token
)

func init() {
	keywords = make(map[string]Token, keywords_end-(keywords_start+1))
	for i := keywords_start + 1; i < keywords_end; i++ {
		keywords[tokens[i]] = i
	}

	literals = make(map[string]Token, lit_end-(lit_start+1))
	for i := lit_start + 1; i < lit_end; i++ {
		literals[tokens[i]] = i
	}

	types = make(map[string]Token, types_end-(types_start+1))
	for i := types_start + 1; i < types_end; i++ {
		types[tokens[i]] = i
	}
}

func (t Token) IsType() bool {
	return t > types_start && t < types_end
}

func (t Token) IsLiteral() bool {
	return t > types_start && t < types_end
}

func (t Token) String() string {
	return tokens[t]
}

func (t Token) IsKeyword() bool {
	return t > keywords_start && t < keywords_end
}

func (t Token) IsOperator() bool {
	return t > keywords_start && t < keywords_end
}

func Lookup(lit string) Token {
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
