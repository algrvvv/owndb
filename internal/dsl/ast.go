package dsl

import "fmt"

type Statement interface {
	String() string
}

type Expression interface {
	String() string
}

// как я понял это числа или любое другое значение для чего либо
type Literal[T any] struct {
	Value T
}

func (l *Literal[T]) String() string { return fmt.Sprintf("{%T: %v}", l.Value, l.Value) }

// как идентификаторы, то есть, в нашем случае,
// это названия переменных.
type Identifier struct {
	Name string
}

func (i *Identifier) String() string { return i.Name }

type BinaryExpr struct {
	Left  Expression
	Op    string
	Right Expression
}

func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(1 - %q; 2 - %q) - %q", b.Left.String(), b.Right.String(), b.Op)
}

type IfStatement struct{}

// SetStatement реализация состояния для SET команды
type SetStatement struct {
	Name  string
	Type  TokenType
	Value Expression
}

func (s *SetStatement) String() string {
	return fmt.Sprintf("SET (type: %s) (name: %s) = (value: %v)", s.Type, s.Name, s.Value.String())
}

type GetStatement struct {
	Name string
}

func (g *GetStatement) String() string {
	return fmt.Sprintf("GET (name: %s)", g.Name)
}

type ShowStatement struct {
	Limit int
}

func (s *ShowStatement) String() string {
	return fmt.Sprintf("SHOW (limit: %d)", s.Limit)
}

type KeysStatement struct {
	Limit int
}

func (k *KeysStatement) String() string {
	return fmt.Sprintf("KEYS (limit: %d)", k.Limit)
}

type RemoveStatement struct {
	Name string
}

func (r *RemoveStatement) String() string {
	return fmt.Sprintf("REMOVE (KEY: %s)", r.Name)
}

type SaveStatement struct{}

func (s *SaveStatement) String() string {
	return fmt.Sprintf("SAVE")
}
