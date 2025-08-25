package gogetosm

import "github.com/op/go-logging"

var log = logging.MustGetLogger("gogetosm")

func SetLogLevel(level logging.Level) {
	logging.SetLevel(level, "gogetosm")
}
