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
	defaultUnit        = 100
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
	}, ";")

	s.Canvas.Path(d, styles)
}

func (s *Schematic) battery(x1, y1, x2, y2 int) {
	const (
		termGap       = defaultUnit / 6
		negTermHeight = defaultUnit / 4
		posTermHeight = defaultUnit / 2
	)

	// s.Canvas.Circle(x1, y1, 10)
	// s.Canvas.Circle(x2, y2, 10)

	angle := math.Atan2(float64(y2-y1), float64(x2-x1))
	distance := int(math.Round(math.Hypot(float64(x2-x1), float64(y2-y1))))
	negTermX := x1 + distance/2 - termGap/2
	posTermX := negTermX + termGap

	s.Canvas.Gtransform(fmt.Sprintf(
		"rotate(%d, %d, %d)",
		int(angle*180/math.Pi), x1, y1,
	))
	s.path(fmt.Sprintf(`
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
	))
	s.Canvas.Gend()
}

func (s *Schematic) AddElement(e *parsing.Element) {

	switch e.Type {
	case "battery":
		defer func(prevX int, prevY int, X *int, Y *int) {
			s.battery(prevX, prevY, *X, *Y)
		}(s.X, s.Y, &s.X, &s.Y)
	case "line":
		defer func(prevX int, prevY int, X *int, Y *int) {
			s.Line(prevX, prevY, *X, *Y)
		}(s.X, s.Y, &s.X, &s.Y)
	}

	for _, action := range e.Actions {

		value := action.Value
		if value == 0 {
			value = defaultUnit
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
}

func NewSchematic(width, height int) *Schematic {
	buffer := &bytes.Buffer{}
	canvas := svg.New(buffer)
	canvas.Start(width, height)
	return &Schematic{
		X:      width / 2,
		Y:      height / 2,
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
