package drawing

import (
	"bytes"
	"fmt"
	"math"
	"os"

	"github.com/sermuns/schemgo/parsing"
)

const (
	unitLength         = 100
	defaultLength      = 1 * unitLength
	defaultStrokeWidth = 5
)

type point struct {
	x, y float64
}

type command struct {
	penDown bool
	pos     point
}

type Path struct {
	commands   []command
	attributes string
}

// 1. translate everything so no children have negative positions
// 2. clamp so canvas is just as big as furthest child
// func (c *Schematic) Normalise() {
// }

type Resistor struct {
	start, end point
}

type Schematic struct {
	X, Y  float64
	paths []*Path
}

// `input` is grouped into 3 and parsed as pathCommand
func createPath(input ...any) Path {
	if len(input)%3 != 0 {
		panic("createPath: input must be a multiple of 3")
	}

	newPath := Path{}
	newCommand := command{pos: point{0, 0}, penDown: true}
	for i, part := range input {
		partType := i % 3
		switch partType {
		case 0:
			newCommand.penDown = part.(bool)
		case 1:
			newCommand.pos.x = part.(float64)
		case 2:
			newCommand.pos.y = part.(float64)
			newPath.commands = append(newPath.commands, newCommand)
		}
	}
	return newPath
}

func (s *Schematic) addPath(path Path) {
	s.paths = append(s.paths, &path)
}

func (this *point) distanceTo(other point) float64 {
	return math.Hypot(other.x-this.x, other.y-this.y)
}

func (s *Schematic) resistor(p1, p2 point) {
	const (
		height = defaultLength / 4
		width  = defaultLength / 2
	)

	distance := p1.distanceTo(p2)
	angleDeg := int(180 / math.Pi * math.Atan2(p2.y-p1.y, p2.x-p1.x))

	path := createPath(
		false, p1.x, p1.y,
		true, p1.x+distance/2-width/2, p1.y,
		true, p1.x+distance/2-width/2, p1.y+height/2,
		true, p1.x+distance/2+width/2, p1.y+height/2,
		true, p1.x+distance/2+width/2, p1.y-height/2,
		true, p1.x+distance/2-width/2, p1.y-height/2,
		true, p1.x+distance/2-width/2, p1.y,
		false, p1.x+distance/2+width/2, p1.y,
		true, p1.x+distance, p1.y,
	)
	path.attributes = fmt.Sprintf(
		` transform="rotate(%d, %d, %d)"`,
		angleDeg, int(p1.x), int(p2.x),
	)
	s.addPath(path)
}

// func (s *Schematic) battery(x1, y1, x2, y2 int) {
// 	const (
// 		termGap       = defaultLength / 6
// 		negTermHeight = defaultLength / 4
// 		posTermHeight = defaultLength / 2
// 	)
//
// 	distance := getIntegerDistance(x1, y1, x2, y2)
// 	negTermX := x1 + distance/2 - termGap/2
// 	posTermX := negTermX + termGap
//
// 	path := fmt.Sprintf(`
// 		M %d %d
// 		L %d %d
// 		M %d %d
// 		L %d %d
// 		M %d %d
// 		L %d %d
// 		M %d %d
// 		L %d %d
// 		`,
// 		x1, y1,
// 		negTermX, y1,
// 		negTermX, y1-negTermHeight/2,
// 		negTermX, y1+negTermHeight/2,
// 		posTermX, y1-posTermHeight/2,
// 		posTermX, y1+posTermHeight/2,
// 		posTermX, y1,
// 		x1+distance, y1,
// 	)
//
// 	s.biNodal(x1, y1, x2, y2, path)
// }

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

	p1 := point{prevX, prevY}
	p2 := point{s.X, s.Y}

	switch e.Type {
	case "resistor":
		s.resistor(p1, p2)
		// case "battery":
		// 	s.battery(prevX, prevY, s.X, s.Y)
		// case "line":
		// 	s.Line(prevX, prevY, s.X, s.Y)
	}
}

func NewSchematic(width, height int) *Schematic {
	return &Schematic{
		X:     defaultLength / 2,
		Y:     defaultLength * 3,
		paths: []*Path{},
	}
}

// func (s *Schematic) Line(x1, y1, x2, y2 int, styles ...string) *Schematic {
// 	styles = append(
// 		styles,
// 		"stroke:"+s.Color,
// 		"stroke-linecap: square",
// 		fmt.Sprintf("stroke-width: %d", defaultStrokeWidth),
// 	)
// 	s.Canvas.Line(x1, y1, x2, y2, strings.Join(styles, ";"))
// 	return s
// }

func (s *Schematic) Translate(dx, dy float64) *Schematic {
	s.X += dx
	s.Y += dy
	return s
}

// create string svg representation.
// need to create root <svg> tag
// and convert all paths to <path d=...> tags
func (s *Schematic) End(outFilePath string) {
	var buf bytes.Buffer

	buf.WriteString(`<svg width='500' height='500'>`)

	for _, path := range s.paths {
		buf.WriteString(fmt.Sprintf(`<path `))
		buf.WriteString(`d="`)
		for _, command := range path.commands {
			if command.penDown {
				buf.WriteString("L ")
			} else {
				buf.WriteString("M ")
			}
			buf.WriteString(
				fmt.Sprintf("%d %d ",
					int(command.pos.x),
					int(command.pos.y),
				),
			)
		}
		buf.WriteString(`"`)
		buf.WriteString(path.attributes)
		buf.WriteString(` style="stroke:black; stroke-width:5; fill:none;"></path>`)
	}

	buf.WriteString("</svg>")

	os.WriteFile(outFilePath, buf.Bytes(), 0644)
}
