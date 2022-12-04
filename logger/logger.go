package logger

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once

var log zerolog.Logger

// Get initializes the logger and returns it.
// It is safe to call this function multiple times.
// It will only initialize the logger once.
// use levelFn to set the log level.
// Log level can be set to one of the following:
// disabled 7
// nolevel 6
// panic  5
// fatal  4
// error  3
// warn  2
// info  1
// debugl 0
// tracel -1
func Get(levelFn ...func() int) zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano
		zerolog.TimestampFunc = func() time.Time {
			return time.Now().UTC()
		}

		var logLevel int
		switch {
		case len(levelFn) > 0:
			logLevel = levelFn[0]()
		default:
			logLevel = 1 // default to info
		}

		log = zerolog.New(os.Stderr).
			Level(IntToLevel(logLevel)).
			With().
			Timestamp().
			Logger()
	})
	return log
}

func IntToLevel(level int) zerolog.Level {
	switch level {
	case 7:
		return zerolog.Disabled
	case 6:
		return zerolog.NoLevel
	case 5:
		return zerolog.PanicLevel
	case 4:
		return zerolog.FatalLevel
	case 3:
		return zerolog.ErrorLevel
	case 2:
		return zerolog.WarnLevel
	case 1:
		return zerolog.InfoLevel
	case 0:
		return zerolog.DebugLevel
	case -1:
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}
