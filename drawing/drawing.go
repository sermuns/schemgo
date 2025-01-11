package drawing

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/svgo"
	"github.com/sermuns/schemgo/parsing"
)

const (
	unitLength         = 100
	defaultLength      = 1 * unitLength
	defaultStrokeWidth = 5
)

type Schematic struct {
	X, Y   int
	Color  string
	Canvas *svg.SVG
	Buffer *bytes.Buffer
}

func (s *Schematic) path(d string) {
	styles := strings.Join([]string{
		"stroke:" + s.Color,
		"stroke-width:" + strconv.Itoa(defaultStrokeWidth),
		"fill: none",
	}, ";")

	s.Canvas.Path(d, styles)
}

// Element that connects two nodes.
// `path` must be in horizontal, left-to-right layout,
// this function takes care of rotation.
func (s *Schematic) biNodal(x1, y1, x2, y2 int, path string) {
	angle := math.Atan2(
		float64(y2-y1),
		float64(x2-x1),
	)

	s.Canvas.Gtransform(fmt.Sprintf(
		"rotate(%d, %d, %d)",
		int(angle*180/math.Pi), x1, y1,
	))
	s.path(path)
	s.Canvas.Gend()
}

func getIntegerDistance(x1, y1, x2, y2 int) int {
	return int(math.Round(math.Hypot(float64(y2-y1), float64(x2-x1))))
}

func (s *Schematic) resistor(x1, y1, x2, y2 int) {
	const (
		height = defaultLength / 4
		width  = defaultLength / 2
	)

	distance := getIntegerDistance(x1, y1, x2, y2)

	path := fmt.Sprintf(`
		M %d %d
		L %d %d
		L %d %d
		L %d %d
		L %d %d
		L %d %d
		L %d %d
		M %d %d
		L %d %d
		`,
		x1, y1,
		x1+distance/2-width/2, y1,
		x1+distance/2-width/2, y1+height/2,
		x1+distance/2+width/2, y1+height/2,
		x1+distance/2+width/2, y1-height/2,
		x1+distance/2-width/2, y1-height/2,
		x1+distance/2-width/2, y1,
		x1+distance/2+width/2, y1,
		x1+distance, y1,
	)

	s.biNodal(x1, y1, x2, y2, path)

}

func (s *Schematic) battery(x1, y1, x2, y2 int) {
	const (
		termGap       = defaultLength / 6
		negTermHeight = defaultLength / 4
		posTermHeight = defaultLength / 2
	)

	distance := getIntegerDistance(x1, y1, x2, y2)
	negTermX := x1 + distance/2 - termGap/2
	posTermX := negTermX + termGap

	path := fmt.Sprintf(`
		M %d %d
		L %d %d
		M %d %d
		L %d %d
		M %d %d
		L %d %d
		M %d %d
		L %d %d
		`,
		x1, y1,
		negTermX, y1,
		negTermX, y1-negTermHeight/2,
		negTermX, y1+negTermHeight/2,
		posTermX, y1-posTermHeight/2,
		posTermX, y1+posTermHeight/2,
		posTermX, y1,
		x1+distance, y1,
	)

	s.biNodal(x1, y1, x2, y2, path)
}

func (s *Schematic) AddElement(e *parsing.Element) {

	prevX, prevY := s.X, s.Y

	for _, action := range e.Actions {

		value := action.Units * unitLength
		if value == 0 {
			value = defaultLength
		}

		switch action.Type {
		case "right":
			s.Translate(value, 0)
		case "up":
			s.Translate(0, -value)
		case "left":
			s.Translate(-value, 0)
		case "down":
			s.Translate(0, value)
		}
	}

	switch e.Type {
	case "resistor":
		s.resistor(prevX, prevY, s.X, s.Y)
	case "battery":
		s.battery(prevX, prevY, s.X, s.Y)
	case "line":
		s.Line(prevX, prevY, s.X, s.Y)
	}
}

func NewSchematic(width, height int) *Schematic {
	buffer := &bytes.Buffer{}
	canvas := svg.New(buffer)
	canvas.Start(width, height)
	return &Schematic{
		X:      defaultLength / 2,
		Y:      defaultLength * 3,
		Color:  "black",
		Canvas: canvas,
		Buffer: buffer,
	}
}

// TODO: optimize this
func (s *Schematic) ChangeCanvasSize(w, h int) *Schematic {
	newSvgTag := fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">`, w, h)
	modifiedLines := strings.Split(s.Buffer.String(), "\n")
	for i, line := range modifiedLines {
		if strings.HasPrefix(line, "<svg") {
			fmt.Println("mdoifying!")
			modifiedLines[i] = newSvgTag
			s.Buffer.Reset()
			s.Buffer.WriteString(strings.Join(modifiedLines, "\n"))
			return s
		}
	}
	panic("Couldn't find <svg> tag!")
}

func (s *Schematic) ChangeBrushColor(color string) *Schematic {
	s.Color = color
	return s
}

func (s *Schematic) Circle(radius int) *Schematic {
	s.Canvas.Circle(s.X, s.Y, radius, "fill:"+s.Color)
	return s
}

func (s *Schematic) Line(x1, y1, x2, y2 int, styles ...string) *Schematic {
	styles = append(
		styles,
		"stroke:"+s.Color,
		"stroke-linecap: square",
		fmt.Sprintf("stroke-width: %d", defaultStrokeWidth),
	)
	s.Canvas.Line(x1, y1, x2, y2, strings.Join(styles, ";"))
	return s
}

func (s *Schematic) Translate(dx int, dy int) *Schematic {
	s.X += dx
	s.Y += dy
	return s
}

func (s *Schematic) End(outputFile string) error {
	s.Canvas.End()
	return os.WriteFile(outputFile, s.Buffer.Bytes(), 0644)
}
