package main

import (
	"os"

	"github.com/algrvvv/owndb/internal/config"
	"github.com/algrvvv/owndb/internal/exec"
	"github.com/algrvvv/owndb/internal/logger"
	"github.com/algrvvv/owndb/internal/repl"
)

func main() {
	conf := config.MustLoad()
	log := logger.MustInit(conf.LogFile, conf.DebugMode)

	executor, err := exec.NewTcpExecutor(conf.Port, log)
	if err != nil {
		panic(err)
	}

	repl := repl.NewREPLInstance(executor, log)
	err = repl.Scan(os.Stdin)
	if err != nil {
		panic(err)
	}
}
