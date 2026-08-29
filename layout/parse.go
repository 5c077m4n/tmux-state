package layout

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrParseLayout     = errors.New("could not parse layout")
	ErrParseOutOfBound = errors.New("pointer is out of the input's bounds")
)

func isDecimalDigit(r rune) bool { return r >= '0' && r <= '9' }
func isHexDigit(r rune) bool     { return r >= '0' && r <= '9' || r >= 'a' && r <= 'f' }

type parser struct {
	input   string
	pointer int
}

func (p *parser) current() rune {
	if p.pointer >= len(p.input) {
		return 0
	}
	return rune(p.input[p.pointer])
}
func (p *parser) peek() rune {
	if p.pointer >= len(p.input)-1 {
		return 0
	}
	return rune(p.input[p.pointer+1])
}
func (p *parser) next() error {
	if p.pointer < len(p.input)-1 {
		p.pointer++
		return nil
	}
	return ErrParseOutOfBound
}
func (p *parser) expect(expected ...rune) error {
	if !slices.Contains(expected, p.current()) {
		expectedStrings := []string{}
		for _, c := range expected {
			expectedStrings = append(expectedStrings, "`"+string(c)+"`")
		}

		return fmt.Errorf(
			"expected %s but got `%c`",
			strings.Join(expectedStrings, " or "),
			p.current(),
		)
	}

	return p.next()
}

func (p *parser) parseMaybeChecksum() (string, bool) {
	maybeChecksum := p.input[p.pointer : p.pointer+4]
	for _, c := range maybeChecksum {
		if !isHexDigit(c) {
			return "", false
		}
	}
	if p.input[p.pointer+len(maybeChecksum)] != ',' {
		return "", false
	}

	p.pointer += len(maybeChecksum + ",")
	return maybeChecksum, true
}
func (p *parser) parseNumber() (int, error) {
	var numberBuilder strings.Builder
	for i := p.pointer; i < len(p.input); i++ {
		if c := rune(p.input[i]); isDecimalDigit(c) {
			numberBuilder.WriteRune(c)

			if err := p.next(); err != nil {
				return 0, err
			}
		} else {
			break
		}
	}

	return strconv.Atoi(numberBuilder.String())
}
func (p *parser) parseDimentions(layout *Layout) (*Layout, error) {
	width, err := p.parseNumber()
	if err != nil {
		return nil, err
	}
	layout.Width = width

	if err := p.expect('x'); err != nil {
		return nil, err
	}

	height, err := p.parseNumber()
	if err != nil {
		return nil, err
	}
	layout.Height = height

	if err := p.expect(','); err != nil {
		return nil, err
	}

	x, err := p.parseNumber()
	if err != nil {
		return nil, err
	}
	layout.X = x

	if err := p.expect(','); err != nil {
		return nil, err
	}

	y, err := p.parseNumber()
	if err != nil {
		return nil, err
	}
	layout.Y = y

	return layout, nil
}
func (p *parser) parseInnerLayout() (*Layout, error) {
	layout := &Layout{}

	layout, err := p.parseDimentions(layout)
	if err != nil {
		return nil, err
	}

	if p.current() == ',' {
		if err := p.next(); err != nil {
			return nil, err
		}

		paneID, err := p.parseNumber()
		if err != nil {
			return nil, err
		}
		layout.PaneID = &paneID
	} else if p.current() == '[' || p.current() == '{' {
		children, err := p.parseChildLayouts()
		if err != nil {
			return nil, err
		}
		layout.Children = children
	}

	return layout, nil
}
func (p *parser) parseChildLayouts() ([]*Layout, error) {
	layout := &Layout{}

	if p.current() == '[' {
		layout.Direction = Vertical
	} else if p.current() == '{' {
		layout.Direction = Horizontal
	}

	if err := p.expect('[', '{'); err != nil {
		return nil, err
	}

	layouts := []*Layout{}

	layout, err := p.parseDimentions(layout)
	if err != nil {
		return nil, err
	}

	if err := p.expect(',', ']', '}'); err != nil {
		return nil, err
	}

	if isDecimalDigit(p.current()) {
		paneID, err := p.parseNumber()
		if err != nil {
			return nil, err
		}
		layout.PaneID = &paneID

		layouts = append(layouts, layout)

		for p.current() == ',' {
			if err := p.next(); err != nil {
				return nil, err
			}

			innerLayout, err := p.parseInnerLayout()
			if err != nil {
				return nil, err
			}

			layouts = append(layouts, innerLayout)
		}
	} else if p.current() == '[' || p.current() == '{' {
		childLayouts, err := p.parseChildLayouts()
		if err != nil {
			return nil, err
		}
		layout.Children = childLayouts
	}

	if err := p.expect(']', '}', 0); err != nil && !errors.Is(err, ErrParseOutOfBound) {
		return nil, err
	}
	return layouts, nil
}
func (p *parser) parseOuterLayout() (*Layout, error) {
	layout := &Layout{}
	if maybeChecksum, found := p.parseMaybeChecksum(); found {
		layout.Checksum = maybeChecksum
	}

	layout, err := p.parseDimentions(layout)
	if err != nil {
		return nil, err
	}

	children, err := p.parseChildLayouts()
	if err != nil {
		return nil, err
	}
	layout.Children = children

	return layout, nil
}

func Parse(input string) (*Layout, error) {
	p := parser{input: strings.ToLower(input), pointer: 0}

	layout, err := p.parseOuterLayout()
	if err != nil {
		return nil, errors.Join(
			ErrParseLayout,
			fmt.Errorf("@ %#v", p),
			err,
		)
	}

	return layout, nil
}
