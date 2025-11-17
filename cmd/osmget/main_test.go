// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
package main

import (
	"testing"

	"github.com/EricNeid/go-getosm/internal/verify"
)

func TestDoesFileExists_shouldReturnTrue(t *testing.T) {
	// action
	exits := doesFileExists("../../test/testdata/sample.txt")
	//
	verify.Assert(t, exits, "file wrongly detected as not present")
}

func TestDoesFileExists_emptyFile_shouldReturnFalse(t *testing.T) {
	// action
	exits := doesFileExists("../../test/testdata/sample-empty.txt")
	//
	verify.Assert(t, !exits, "file wrongly detected as not empty")
}

func TestDoesFileExists_noFile_shouldReturnFalse(t *testing.T) {
	// action
	exits := doesFileExists("test/testdata/no-valid-file.txt")
	//
	verify.Assert(t, !exits, "file wrongly detected as present")
}
