package layout

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var layouts = [...]string{
	"159x48,0,0[159x24,0,0,1,159x23,0,25{79x23,0,25,2,79x23,80,25,3}]",
	"239x58,0,0[239x47,0,0,1,239x10,0,48,2]",
	"239x58,0,0[239x44,0,0,3,239x13,0,45{143x13,0,45,4,95x13,144,45,6}]",
}
var checksums = [...]string{"bb62", "34b0", "b0a8"}

func TestCRC16(t *testing.T) {
	for index, layout := range layouts {
		t.Run("Checksum for "+layout, func(t *testing.T) {
			result := calculateTmuxCRC16(layout)
			assert.Equal(t, checksums[index], fmt.Sprintf("%04x", result))
		})
	}
}
