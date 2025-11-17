// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
package gogetosm

import (
	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("gogetosm")

func SetLogLevel(level logging.Level) {
	logging.SetLevel(level, "gogetosm")
}
