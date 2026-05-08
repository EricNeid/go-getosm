// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
package gogetosm

import (
	"github.com/op/go-logging"
)

// Log is the logger for this application.
var Log = logging.MustGetLogger("gogetosm")

// SetLogLevel sets the log level for the application logger.
func SetLogLevel(level logging.Level) {
	logging.SetLevel(level, "gogetosm")
}
