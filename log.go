// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
package gogetosm

import (
	"log/slog"
)

// Log is the logger for this application.
// It uses the default format, eg:
//
//	2022/11/08 15:28:26 INFO hello count=3
var Log = slog.Default()

// SetLogLevel sets the log level for the application logger.
func SetLogLevel(level slog.Level) {
	slog.SetLogLoggerLevel(level)
}
