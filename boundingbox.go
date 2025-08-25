package gogetosm

import (
	"errors"
	"strconv"
	"strings"
)

var ErrorInvalidBB = errors.New("invalid bounding box given")

type BoundingBox struct {
	w, s, e, n float64
}

func ReadBoundingBox(bbString string, tiles int) (bbs []BoundingBox, err error) {
	bbStrParts := strings.Split(bbString, ",")
	if len(bbStrParts) != 4 {
		log.Errorf("invalid bounding box given: expecting w,s,e,n")
		return bbs, ErrorInvalidBB
	}
	w, err := strconv.ParseFloat(bbStrParts[0], 64)
	if err != nil {
		log.Errorf("could not parse west %s\n", bbStrParts[0])
		return bbs, ErrorInvalidBB
	}
	s, err := strconv.ParseFloat(bbStrParts[1], 64)
	if err != nil {
		log.Errorf("could not parse south %s\n", bbStrParts[1])
		return bbs, ErrorInvalidBB
	}
	e, err := strconv.ParseFloat(bbStrParts[2], 64)
	if err != nil {
		log.Errorf("could not parse east %s\n", bbStrParts[2])
		return bbs, ErrorInvalidBB
	}
	n, err := strconv.ParseFloat(bbStrParts[3], 64)
	if err != nil {
		log.Errorf("could not parse north %s\n", bbStrParts[3])
		return bbs, ErrorInvalidBB
	}

	if tiles == 1 {
		return []BoundingBox{{w, s, e, n}}, nil
	}

	slidingWest := w
	for i := 0; i < tiles; i++ {
		e = slidingWest + (e-w)/float64(tiles)
		bbs = append(bbs, BoundingBox{
			w: slidingWest,
			s: s,
			e: e,
			n: n,
		})
		slidingWest = e
	}
	return bbs, nil
}
