package main

import (
	"log"
	"strconv"
	"strings"
)

type boundingBox struct {
	w, s, e, n float64
}

func readBoundingBox(bbString string) (bb boundingBox, err error) {
	bbStrParts := strings.Split(bbString, ",")
	if len(bbStrParts) != 4 {
		log.Println("invalid bounding box given: expecting w,s,e,n")
		return bb, errorInvalidBB
	}
	w, err := strconv.ParseFloat(bbStrParts[0], 64)
	if err != nil {
		log.Printf("could not parse west %s\n", bbStrParts[0])
		return bb, errorInvalidBB
	}
	s, err := strconv.ParseFloat(bbStrParts[1], 64)
	if err != nil {
		log.Printf("could not parse south %s\n", bbStrParts[1])
		return bb, errorInvalidBB
	}
	e, err := strconv.ParseFloat(bbStrParts[2], 64)
	if err != nil {
		log.Printf("could not parse east %s\n", bbStrParts[2])
		return bb, errorInvalidBB
	}
	n, err := strconv.ParseFloat(bbStrParts[3], 64)
	if err != nil {
		log.Printf("could not parse north %s\n", bbStrParts[3])
		return bb, errorInvalidBB
	}
	return boundingBox{w, s, e, n}, nil
}
