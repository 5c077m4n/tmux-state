package layout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var layoutStringToStruct = map[string]*Layout{
	"34b0,239x58,0,0[239x47,0,0,1,239x10,0,48,2]":                             {},
	"b0a8,239x58,0,0[239x44,0,0,3,239x13,0,45{143x13,0,45,4,95x13,144,45,6}]": {},
	"2130,239x58,0,0[239x37,0,0,1,239x20,0,38,2]":                             {},
	"7068,239x58,0,0[239x44,0,0,0,239x13,0,45{143x13,0,45,4,95x13,144,45,5}]": {},
}

func TestLayoutParse(t *testing.T) {
	for layoutString, expected := range layoutStringToStruct {
		t.Run(layoutString+" should be parsed to "+expected.String(), func(t *testing.T) {
			parsed, err := NewParser(layoutString).Parse()

			if assert.Nil(t, err) {
				assert.Equal(t, expected, parsed)
			}
		})
	}
}
