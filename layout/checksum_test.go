package layout

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var layoutToChecksum = map[string]string{
	"239x58,0,0[239x47,0,0,1,239x10,0,48,2]":                             "34b0",
	"239x58,0,0[239x44,0,0,3,239x13,0,45{143x13,0,45,4,95x13,144,45,6}]": "b0a8",
	"239x58,0,0[239x37,0,0,1,239x20,0,38,2]":                             "2130",
	"239x58,0,0[239x44,0,0,0,239x13,0,45{143x13,0,45,4,95x13,144,45,5}]": "7068",
}

func TestLayoutChecksum(t *testing.T) {
	for layout, expectedChecksum := range layoutToChecksum {
		t.Run("Checksum for "+layout+" should equal "+expectedChecksum, func(t *testing.T) {
			result := calculateTmuxChecksum(layout)
			assert.Equal(t, expectedChecksum, fmt.Sprintf("%04x", result))
		})
	}
}
