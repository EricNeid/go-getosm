// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
// Inspired by SUMO's osmGet.py (no code copied).
package gogetosm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrorInvalidBB = errors.New("invalid bounding box given")

type BoundingBox struct {
	West, South, East, North float64
}

type TileMode string

const (
	TileVertical   = "vertical"
	TileHorizontal = "horizontal"
	TilegGrid      = "grid"
)

func ParseTileMode(str string) (TileMode, error) {
	switch TileMode(str) {
	case TileVertical, TileHorizontal, TilegGrid:
		return TileMode(str), nil
	default:
		Log.Errorf("could not parse tile mode %s\n", str)
		return "", fmt.Errorf("unknown TileMode: %s", str)
	}
}

func ReadBoundingBox(bbString string, tiles int, tileMode TileMode) (bbs []BoundingBox, err error) {
	bbStrParts := strings.Split(bbString, ",")
	if len(bbStrParts) != 4 {
		Log.Errorf("invalid bounding box given: expecting w,s,e,n")
		return bbs, ErrorInvalidBB
	}
	w, err := strconv.ParseFloat(bbStrParts[0], 64)
	if err != nil {
		Log.Errorf("could not parse west %s\n", bbStrParts[0])
		return bbs, ErrorInvalidBB
	}
	s, err := strconv.ParseFloat(bbStrParts[1], 64)
	if err != nil {
		Log.Errorf("could not parse south %s\n", bbStrParts[1])
		return bbs, ErrorInvalidBB
	}
	e, err := strconv.ParseFloat(bbStrParts[2], 64)
	if err != nil {
		Log.Errorf("could not parse east %s\n", bbStrParts[2])
		return bbs, ErrorInvalidBB
	}
	n, err := strconv.ParseFloat(bbStrParts[3], 64)
	if err != nil {
		Log.Errorf("could not parse north %s\n", bbStrParts[3])
		return bbs, ErrorInvalidBB
	}

	if tiles == 1 {
		return []BoundingBox{{w, s, e, n}}, nil
	}

	switch tileMode {
	case TileVertical:
		tileWidth := (e - w) / float64(tiles)
		slidingWest := w
		var slidingEast float64
		for i := 0; i < tiles; i++ {
			slidingEast = slidingWest + tileWidth
			bbs = append(bbs, BoundingBox{
				West:  slidingWest,
				South: s,
				East:  slidingEast,
				North: n,
			})
			slidingWest = slidingEast
		}
	case TileHorizontal:
		tileWidth := (n - s) / float64(tiles)
		slidingSouth := s
		var slidingNorth float64
		for i := 0; i < tiles; i++ {
			slidingNorth = slidingSouth + tileWidth
			bbs = append(bbs, BoundingBox{
				West:  w,
				South: slidingSouth,
				East:  e,
				North: slidingNorth,
			})
			slidingSouth = slidingNorth
		}
	default:
		return nil, errors.New("unsupported tileing operation")
	}

	return bbs, nil
}
