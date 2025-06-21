package dsl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/algrvvv/owndb/internal/storage"
	"github.com/algrvvv/owndb/internal/storage/snapshot"
)

type Interpreter struct {
	snap      snapshot.Snapshotter
	storage   storage.Storage
	debugMode bool
}

func NewInterpreter(snap snapshot.Snapshotter, storage storage.Storage) *Interpreter {
	return &Interpreter{
		storage: storage,
		snap:    snap,
	}
}

func (i *Interpreter) ExecStatement(stmt Statement, debug bool) (any, error) {
	start := time.Now()
	i.debugMode = debug

	i.debugPrintf("interpreter debug mode: %v\n", i.debugMode)
	i.debugPrintf("stmt: %s\n", stmt.String())

	switch s := stmt.(type) {
	case *SetStatement:
		valExpr, err := i.evalExpr(s.Value)
		if err != nil {
			return nil, err
		}

		err = i.storage.Set(s.Name, valExpr)
		if err == nil {
			return i.okTime(start), nil
		}

		return nil, err
	case *GetStatement:
		data, ok := i.storage.Get(s.Name)
		if !ok {
			return nil, fmt.Errorf("undefined key: %q", s.Name)
		}

		return data, nil
	case *ShowStatement:
		var count int
		builder := strings.Builder{}
		data := i.storage.GetAll()

		for k, v := range data {
			if count == s.Limit {
				break
			}

			count++
			builder.WriteString(fmt.Sprintf("%q = %v; ", k, v))
		}

		return builder.String(), nil
	case *KeysStatement:
		builder := strings.Builder{}
		keys := i.storage.Keys()

		for i, k := range keys {
			if i == s.Limit {
				break
			}

			builder.WriteString(fmt.Sprintf("%q; ", k))
		}

		return builder.String(), nil

	case *RemoveStatement:
		err := i.storage.Remove(s.Name)
		if err != nil {
			return nil, err
		}

		return i.okTime(start), nil
	case *SaveStatement:
		m := i.storage.GetAll()
		_, err := i.snap.Write(m)
		if err != nil {
			return nil, err
		}

		return i.okTime(start), nil
	default:
		return nil, fmt.Errorf("invalid statement")
	}
}

func (i *Interpreter) evalExpr(expr Expression) (any, error) {
	switch e := expr.(type) {
	// вот сука придумал на свою голову
	// использовать здесь дженерики
	case *Literal[string]:
		return e.Value, nil
	case *Literal[int]:
		return e.Value, nil
	case *Literal[float64]:
		return e.Value, nil
	case *Literal[bool]:
		return e.Value, nil

	case *Identifier:
		val, ok := i.storage.Get(e.Name)
		if !ok {
			return nil, fmt.Errorf("undefined key: %s", e.Name)
		}

		return val, nil
	default:
		return nil, errors.New("")
	}
}
