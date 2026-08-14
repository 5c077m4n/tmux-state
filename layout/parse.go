package layout

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrParseLayout = errors.New("could not parse layout")

func isHexDigit(r rune) bool {
	return r >= '0' && r <= '9' || r >= 'a' && r <= 'f'
}

type parser struct {
	input        string
	pointerStart int
	pointerEnd   int
}

func (p *parser) setPointers(start int) {
	p.pointerStart = start
	p.pointerEnd = p.pointerStart + 1
}

func (p *parser) getSlice() string {
	if len(p.input) < p.pointerStart || len(p.input) < p.pointerEnd {
		return ""
	}
	return p.input[p.pointerStart:p.pointerEnd]
}

func (p *parser) eatUntilNextSection() string {
	for i := p.pointerStart; i < len(p.input) || p.input[i] == ',' || p.input[i] == '[' || p.input[i] == '{'; i++ {
		p.pointerStart++
	}
	p.pointerStart++
	p.pointerEnd = p.pointerStart + 1

	return p.getSlice()
}

func (p *parser) getChecksum() (string, bool) {
	p.setPointers(0)

	maybeChecksum := p.eatUntilNextSection()
	if len(maybeChecksum) != 4 {
		return "", false
	}
	for _, char := range maybeChecksum {
		if !isHexDigit(char) {
			p.setPointers(0)
			return "", false
		}
	}

	if p.input[5] != ',' {
		p.setPointers(0)
		return "", false
	}
	p.setPointers(5)

	return maybeChecksum, true
}

func (p *parser) getXYCoords() (int, int, error) {
	maybeXYCoords := p.eatUntilNextSection()
	parsedXYCoords := strings.Split(maybeXYCoords, "x")
	if len(parsedXYCoords) != 2 {
		return 0, 0, fmt.Errorf("could not parse %s as coordinates", maybeXYCoords)
	}

	x, err := strconv.ParseInt(parsedXYCoords[0], 10, 32)
	if err != nil {
		return 0, 0, err
	}

	y, err := strconv.ParseInt(parsedXYCoords[1], 10, 32)
	if err != nil {
		return 0, 0, err
	}

	return int(x), int(y), nil

}

func (p *parser) Parse() (*Layout, error) {
	layout := &Layout{}

	if maybeChecksum, found := p.getChecksum(); found {
		layout.Checksum = maybeChecksum
	}

	x, y, err := p.getXYCoords()
	if err != nil {
		return nil, errors.Join(ErrParseLayout, err)
	}
	layout.X = x
	layout.Y = y

	heightString := p.eatUntilNextSection()
	height, err := strconv.ParseInt(heightString, 10, 32)
	if err != nil {
		return nil, errors.Join(ErrParseLayout, err)
	}
	layout.Height = int(height)

	widthString := p.eatUntilNextSection()
	width, err := strconv.ParseInt(widthString, 10, 32)
	if err != nil {
		return nil, errors.Join(ErrParseLayout, err)
	}
	layout.Width = int(width)

	if p.input[p.pointerStart] == '[' || p.input[p.pointerStart] == '{' {
		parens := []byte{}

		for i := p.pointerStart; i < len(p.input); i++ {
			currentParens := p.input[i]
			if currentParens == '[' || currentParens == '{' {
				parens = append(parens, currentParens)
			}
			if currentParens == '}' || currentParens == ']' {
				if parens[len(parens)-1] != currentParens {
					return nil, ErrParseLayout
				}
				parens = parens[:len(parens)-1]
			}
			p.pointerEnd++
		}
	}

	return layout, nil
}

func NewParser(input string) *parser {
	return &parser{input: strings.ToLower(input), pointerStart: 0, pointerEnd: 1}
}
