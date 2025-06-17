package dsl

import (
	"fmt"
	"time"
)

func (i *Interpreter) debugPrint(data any) {
	if i.debugMode {
		fmt.Println(data)
	}
}

func (i *Interpreter) debugPrintf(format string, a ...any) {
	if i.debugMode {
		fmt.Printf(format, a...)
	}
}

func (i *Interpreter) okTime(start time.Time) string {
	dur := time.Since(start)
	return fmt.Sprintf("OK! %v", dur.String())
}
