// Package layout parsing/printing
package layout

import (
	"errors"
	"fmt"
	"strings"
)

type LayoutDirection string

const (
	Vertical   LayoutDirection = "vertical"
	Horizontal LayoutDirection = "horizontal"
)

type Layout struct {
	Checksum  string
	Width     int
	Height    int
	X         int
	Y         int
	PaneID    *int
	Direction LayoutDirection
	Children  []*Layout
}

func (l *Layout) IsPane() bool { return l.PaneID != nil }
func (l *Layout) String() string {
	var buf strings.Builder
	if l.Checksum != "" {
		buf.WriteString(l.Checksum)
		buf.WriteByte(',')
	}
	fmt.Fprintf(&buf, "%dx%d,%d,%d", l.Width, l.Height, l.X, l.Y)
	if l.PaneID != nil {
		fmt.Fprintf(&buf, ",%d", *l.PaneID)
	}
	if len(l.Children) > 0 {
		if l.Direction == Horizontal {
			buf.WriteByte('[')
		} else {
			buf.WriteByte('{')
		}

		for i, child := range l.Children {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(child.String())
		}

		if l.Direction == Horizontal {
			buf.WriteByte(']')
		} else {
			buf.WriteByte('}')
		}

	}
	return buf.String()
}

var ErrParseLayout = errors.New("could not parse layout")

// Example layout string:
// bb62,159x48,0,0[159x24,0,0,1,159x23,0,25{79x23,0,25,2,79x23,80,25,3}]
// 34b0,239x58,0,0[239x47,0,0,1,239x10,0,48,2]
// b0a8,239x58,0,0[239x44,0,0,3,239x13,0,45{143x13,0,45,4,95x13,144,45,6}]

func From(_layout string) (*Layout, error) {
	l := &Layout{}

	return l, nil
}
