package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/algrvvv/owndb/internal/exec"
	"github.com/rs/zerolog"
)

type REPL struct {
	executor exec.Executor
	log      zerolog.Logger
}

func NewREPLInstance(executor exec.Executor, log zerolog.Logger) *REPL {
	return &REPL{
		executor: executor,
		log:      log,
	}
}

func (r *REPL) Scan(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		r.log.Debug().Str("line", line).Msg("got new line from user")

		if strings.TrimSpace(line) == "" {
			print("> ")
			continue
		} else if line == "EXIT" || line == "EXIT;" {
			return fmt.Errorf("exit by user")
		}

		data, err := r.executor.Execute(line)
		if err != nil {
			fmt.Println("FAIL! ", err.Error())
			print("\n> ")
			continue
		}

		if data != nil {
			r.log.Debug().Any("data", data).Msg("got data from executor")
			fmt.Println(data)
		}

		print("\n> ")
	}

	if scanner.Err() != nil {
		return fmt.Errorf("scanner error")
	}

	return nil
}
