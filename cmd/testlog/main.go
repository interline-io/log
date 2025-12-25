package main

import (
	"github.com/interline-io/log"
)

func main() {
	log.Print("This is a plain print (no level, no timestamp)")

	log.Tracef("This is a trace message: %d", 1)
	log.Debugf("This is a debug message: %d", 2)
	log.Infof("This is an info message: %d", 3)
	log.Errorf("This is an error message: %d", 4)

	log.Traceln("Traceln:", "multiple", "args", 123)

	log.Info().Str("key", "value").Int("count", 42).Msg("Structured logging example")
	log.Debug().Str("user", "alice").Bool("active", true).Msg("User status")
	log.Error().Err(nil).Str("op", "test").Msg("Error with nil error")

	log.Print("Done!")
}
