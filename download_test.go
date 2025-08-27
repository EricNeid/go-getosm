package gogetosm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/EricNeid/go-getosm/internal/verify"
)

func TestFormatQuery(t *testing.T) {
	// arrange
	bb := BoundingBox{
		North: 52.35211857272093,
		South: 52.29189255277229,
		West:  13.16848754882812,
		East:  13.27835083007812,
	}
	// action
	result := FormatQuery(bb, 240, 1073741824)
	// verify
	verify.NotNil(t, result, "query is null")
	verify.Assert(t, strings.Contains(result, "n=\"52.35211857272093\""), fmt.Sprintf("substring not found: %s", "n=\"52.35211857272093\""))
	verify.Assert(t, strings.Contains(result, "s=\"52.29189255277229\""), fmt.Sprintf("substring not found: %s", "s=\"52.29189255277229\""))
	verify.Assert(t, strings.Contains(result, "w=\"13.16848754882812\""), fmt.Sprintf("substring not found: %s", "w=\"13.16848754882812\""))
	verify.Assert(t, strings.Contains(result, "e=\"13.27835083007812\""), fmt.Sprintf("substring not found: %s", "e=\"13.27835083007812\""))
}
