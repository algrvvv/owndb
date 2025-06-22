package utils

import (
	"github.com/algrvvv/owndb/internal/dsl"
	"github.com/algrvvv/owndb/internal/exec"
	"github.com/algrvvv/owndb/internal/storage"
	"github.com/algrvvv/owndb/internal/wal"
	"github.com/rs/zerolog"
)

func PrepareData(logger zerolog.Logger, storage storage.Storage, wal *wal.WAL, intr *dsl.Interpreter) error {
	logs, err := wal.Read()
	if err != nil {
		return err
	}

	executor := exec.NewInlineDiscardExecutor(intr)
	for _, log := range logs {
		logger.Debug().Msgf("got log from wal: %q", log)
		res, err := executor.Execute(log)
		if err != nil {
			logger.Debug().Msgf("exec log err: %s", err)
			continue
		}

		logger.Debug().Msgf("exec log res: %v", res)
	}

	err = wal.Clear()
	if err != nil {
		return err
	}

	return nil
}
