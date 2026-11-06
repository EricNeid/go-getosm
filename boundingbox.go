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

// ErrorInvalidBB indicates that the given bounding box string is not valid.
var ErrorInvalidBB = errors.New("invalid bounding box given")

// BoundingBox is a wrapper for geographic bounding box.
type BoundingBox struct {
	West, South, East, North float64
}

// ReadBoundingBox reads the given bounding box string and creates a list of boxed, according to the given number
// of rows and columns.
func ReadBoundingBox(bbString string, tilesX, tilesY int) (bbs []BoundingBox, err error) {
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

	tileWidth := (e - w) / float64(tilesX)
	tileHeight := (n - s) / float64(tilesY)
	for y := range tilesY {
		south := s + float64(y)*tileHeight
		north := south + tileHeight
		for x := range tilesX {
			west := w + float64(x)*tileWidth
			east := west + tileWidth
			bbs = append(bbs, BoundingBox{
				West:  west,
				South: south,
				East:  east,
				North: north,
			})

		}
	}
	return bbs, nil
}

// BoundingBoxesToWkt reads the given bounding box strings and creates a wkt multipolygon.
func BoundingBoxesToWkt(bbStrings []string) string {
	var bbs []BoundingBox
	for _, bbStbbString := range bbStrings {
		bb, _ := ReadBoundingBox(bbStbbString, 1, 1)
		bbs = append(bbs, bb...)
	}
	if len(bbs) == 0 {
		return ""
	}
	var polygons []string
	for _, bb := range bbs {
		// WKT polygon ring for single bb
		// (West South, East South, East North, West North, West South)
		poly := fmt.Sprintf("((%f %f, %f %f, %f %f, %f %f, %f %f))",
			bb.West, bb.South,
			bb.East, bb.South,
			bb.East, bb.North,
			bb.West, bb.North,
			bb.West, bb.South,
		)
		polygons = append(polygons, poly)
	}
	return "MULTIPOLYGON(" + strings.Join(polygons, ", ") + ")"
}
