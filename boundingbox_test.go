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
	verify.Equals(t, BoundingBox{10.4, 50.0, 10.6, 51.0}, bbs[0])
	verify.Equals(t, BoundingBox{10.4, 50.0, 10.8, 51.0}, bbs[1])
}
