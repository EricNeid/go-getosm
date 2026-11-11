// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
package main

import (
	"testing"

	"github.com/EricNeid/go-getosm/internal/verify"
)

func TestDoesFileExist_shouldReturnTrue(t *testing.T) {
	// action
	exists := doesFileExist("../../test/testdata/sample.txt")
	// verify
	verify.Assert(t, exists, "file wrongly detected as not present")
}

func TestDoesFileExist_emptyFile_shouldReturnFalse(t *testing.T) {
	// action
	exists := doesFileExist("../../test/testdata/sample-empty.txt")
	// verify
	verify.Assert(t, !exists, "file wrongly detected as not empty")
}

func TestDoesFileExist_noFile_shouldReturnFalse(t *testing.T) {
	// action
	exists := doesFileExist("test/testdata/no-valid-file.txt")
	// verify
	verify.Assert(t, !exists, "file wrongly detected as present")
}

func TestGetFileName(t *testing.T) {
	// action
	name := getFileName("ger", 1, 1, 1, 4)
	// verify
	verify.Equals(t, "ger_2_4.osm.xml", name)
}

func TestGetFileName_multipleBB_shouldAddBBIndex(t *testing.T) {
	// action
	name := getFileName("ger", 1, 2, 1, 4)
	// verify
	verify.Equals(t, "bb_1_2_ger_2_4.osm.xml", name)
}
