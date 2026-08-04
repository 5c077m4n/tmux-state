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

func (l *Layout) asStringWithoutChecksum() string {
	var buf strings.Builder

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

func (l *Layout) IsPane() bool { return l.PaneID != nil }
func (l *Layout) String() string {
	layoutString := l.asStringWithoutChecksum()
	checksum := calculateTmuxChecksum(layoutString)

	return fmt.Sprintf("%04x,%s", checksum, layoutString)
}

var ErrParseLayout = errors.New("could not parse layout")

func From(_layout string) (*Layout, error) {
	l := &Layout{}

	return l, nil
}
