package logging

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

var GlobalLogger zerolog.Logger

func InitializeLogger() {
	GlobalLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
}
