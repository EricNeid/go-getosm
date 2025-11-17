// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
package gogetosm

import (
	"testing"

	"github.com/EricNeid/go-getosm/internal/verify"
)

func TestReadBoundingBox(t *testing.T) {
	// arrange

	// action
	bbs, err := ReadBoundingBox("10.4,50.0,10.8,51.0", 2)
	// verify
	verify.Ok(t, err)
	verify.Equals(t, 2, len(bbs))
	verify.AlmostEquals(t, 51.0, bbs[0].North)
	verify.AlmostEquals(t, 51.0, bbs[1].North)
	verify.AlmostEquals(t, 50.0, bbs[0].South)
	verify.AlmostEquals(t, 50.0, bbs[1].South)
	verify.AlmostEquals(t, 10.4, bbs[0].West)
	verify.AlmostEquals(t, 10.6, bbs[0].East)
	verify.AlmostEquals(t, 10.6, bbs[1].West)
	verify.AlmostEquals(t, 10.8, bbs[1].East)
}
