// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
package gogetosm

import (
	"testing"

	"github.com/EricNeid/go-getosm/internal/verify"
)

func TestParseTileMode(t *testing.T) {
	// action
	res, err := ParseTileMode("vertical")
	// verify
	verify.Ok(t, err)
	verify.Equals(t, TileVertical, res)
}

func TestParseTileMode_invalidMode(t *testing.T) {
	// action
	_, err := ParseTileMode("invalid")
	// verify
	verify.NotNil(t, err, "")
}

func TestReadBoundingBox_tileVertical(t *testing.T) {
	// action
	bbs, err := ReadBoundingBox("10,50,11,51", 2, TileVertical)
	// verify
	verify.Ok(t, err)
	verify.Equals(t, 2, len(bbs))

	verify.AlmostEquals(t, 10.0, bbs[0].West)
	verify.AlmostEquals(t, 10.5, bbs[0].East)
	verify.AlmostEquals(t, 50.0, bbs[0].South)
	verify.AlmostEquals(t, 51.0, bbs[0].North)

	verify.AlmostEquals(t, 10.5, bbs[1].West)
	verify.AlmostEquals(t, 11.0, bbs[1].East)
	verify.AlmostEquals(t, 50.0, bbs[1].South)
	verify.AlmostEquals(t, 51.0, bbs[1].North)
}

func TestReadBoundingBox_tileHorizontal(t *testing.T) {
	// action
	bbs, err := ReadBoundingBox("10,50,11,51.0", 2, TileHorizontal)
	// verify
	verify.Ok(t, err)
	verify.Equals(t, 2, len(bbs))

	verify.AlmostEquals(t, 10.0, bbs[0].West)
	verify.AlmostEquals(t, 11.0, bbs[0].East)
	verify.AlmostEquals(t, 50.0, bbs[0].South)
	verify.AlmostEquals(t, 50.5, bbs[0].North)

	verify.AlmostEquals(t, 10.0, bbs[1].West)
	verify.AlmostEquals(t, 11.0, bbs[1].East)
	verify.AlmostEquals(t, 50.5, bbs[1].South)
	verify.AlmostEquals(t, 51.0, bbs[1].North)
}
