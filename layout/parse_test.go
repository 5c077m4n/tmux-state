package layout

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var layoutStringToStruct = map[string]*Layout{
	"34b0,239x58,0,0[239x47,0,0,1,239x10,0,48,2]": {
		Checksum:  "34b0",
		Width:     239,
		Height:    58,
		X:         0,
		Y:         0,
		Direction: Vertical,
		Children: []*Layout{
			{Width: 239, Height: 47, X: 0, Y: 0, PaneID: new(1)},
			{Width: 239, Height: 10, X: 0, Y: 48, PaneID: new(2)},
		},
	},
	"b0a8,239x58,0,0[239x44,0,0,3,239x13,0,45{143x13,0,45,4,95x13,144,45,6}]": {
		Checksum:  "b0a8",
		Width:     239,
		Height:    58,
		X:         0,
		Y:         0,
		Direction: Vertical,
		Children: []*Layout{
			{Width: 239, Height: 44, X: 0, Y: 0, PaneID: new(3)},
			{
				Width:     239,
				Height:    13,
				X:         0,
				Y:         45,
				Direction: Horizontal,
				Children: []*Layout{
					{Width: 143, Height: 13, X: 0, Y: 45, PaneID: new(4)},
					{Width: 95, Height: 13, X: 144, Y: 45, PaneID: new(6)},
				},
			},
		},
	},
	"2130,239x58,0,0[239x37,0,0,1,239x20,0,38,2]": {
		Checksum:  "2130",
		Width:     239,
		Height:    58,
		X:         0,
		Y:         0,
		Direction: Vertical,
		Children: []*Layout{
			{Width: 239, Height: 37, X: 0, Y: 1, PaneID: new(1)},
			{Width: 239, Height: 20, X: 0, Y: 38, PaneID: new(2)},
		},
	},
	"7068,239x58,0,0[239x44,0,0,0,239x13,0,45{143x13,0,45,4,95x13,144,45,5}]": {
		Checksum:  "7068",
		Width:     239,
		Height:    58,
		X:         0,
		Y:         0,
		Direction: Vertical,
		Children: []*Layout{
			{Width: 239, Height: 44, X: 0, Y: 0, PaneID: new(0)},
			{
				Width:     239,
				Height:    13,
				X:         0,
				Y:         45,
				Direction: Horizontal,
				Children: []*Layout{
					{Width: 143, Height: 13, X: 0, Y: 45, PaneID: new(4)},
					{Width: 95, Height: 13, X: 144, Y: 45, PaneID: new(5)},
				},
			},
		},
	},
}

func TestLayoutParse(t *testing.T) {
	for layoutString, expected := range layoutStringToStruct {
		t.Run(layoutString+" should be parsed to "+expected.String(), func(t *testing.T) {
			parsed, err := NewParser(layoutString).Parse()

			require.NoError(t, err)
			assert.Equal(t, expected, parsed)
		})
	}
}
