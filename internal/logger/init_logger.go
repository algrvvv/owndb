package logger

import (
	"io"
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

func MustInit(path string) zerolog.Logger {
	buildInfo, _ := debug.ReadBuildInfo()

	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("faile to open log file: " + err.Error())
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	multiWriter := io.MultiWriter(consoleWriter, logFile)

	return zerolog.New(multiWriter).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()
}
