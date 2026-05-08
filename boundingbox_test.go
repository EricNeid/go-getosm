// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
package gogetosm

import (
	"testing"

	"github.com/EricNeid/go-getosm/internal/verify"
)

func TestReadBoundingBox_tileVertical(t *testing.T) {
	// action
	bbs, err := ReadBoundingBox("10,50,11,51", 2, 1)
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
	bbs, err := ReadBoundingBox("10,50,11,51.0", 1, 2)
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

func TestReadBoundingBox_tileGrid(t *testing.T) {
	// action
	bbs, err := ReadBoundingBox("10,50,11,51.0", 2, 2)
	// verify
	verify.Ok(t, err)
	verify.Equals(t, 4, len(bbs))

	verify.AlmostEquals(t, 10.0, bbs[0].West)
	verify.AlmostEquals(t, 10.5, bbs[0].East)
	verify.AlmostEquals(t, 50.0, bbs[0].South)
	verify.AlmostEquals(t, 50.5, bbs[0].North)

	verify.AlmostEquals(t, 10.5, bbs[1].West)
	verify.AlmostEquals(t, 11.0, bbs[1].East)
	verify.AlmostEquals(t, 50.0, bbs[1].South)
	verify.AlmostEquals(t, 50.5, bbs[1].North)

	verify.AlmostEquals(t, 10.0, bbs[2].West)
	verify.AlmostEquals(t, 10.5, bbs[2].East)
	verify.AlmostEquals(t, 50.5, bbs[2].South)
	verify.AlmostEquals(t, 51.0, bbs[2].North)

	verify.AlmostEquals(t, 10.5, bbs[3].West)
	verify.AlmostEquals(t, 11.0, bbs[3].East)
	verify.AlmostEquals(t, 50.5, bbs[3].South)
	verify.AlmostEquals(t, 51.0, bbs[3].North)
}
