package drawing

import (
	"bytes"
	"fmt"
	"math"

	"github.com/sermuns/schemgo/parsing"
)

const (
	unitLength         = 100
	defaultLength      = 1 * unitLength
	defaultStrokeWidth = 5
)

var (
	stack Stack
)

type Point struct {
	x, y float64
}

type command struct {
	letter rune
	pos    Point
}

type Path []command

type Circle struct {
	centerPos Point
	radius    float64
}

type Schematic struct {
	pos     Point
	paths   []*Path
	circles []*Circle
}

// `input`:s elements are grouped into triplets and parsed as a path
// i'm sorry this is so fucked up
func createPath(input ...any) Path {
	if len(input)%3 != 0 {
		panic("createPath: input must be a multiple of 3")
	}

	newPath := Path{}
	newCommand := command{}
	for i, part := range input {
		partType := i % 3
		switch partType {
		case 0:
			newCommand.letter = part.(rune)
		case 1:
			newCommand.pos.x = part.(float64)
		case 2:
			newCommand.pos.y = part.(float64)
			newPath = append(newPath, newCommand)
		}
	}
	return newPath
}

func (s *Schematic) addPath(path Path) {
	s.paths = append(s.paths, &path)
}

func (s *Schematic) addCircle(x, y, radius float64) {
	s.circles = append(s.circles, &Circle{
		centerPos: Point{x, y},
		radius:    radius,
	})
}

func (s *Schematic) addAndPivotPath(start, end Point, path Path) {
	angle := math.Atan2(end.y-start.y, end.x-start.x)
	path.pivotAround(start, angle)
	s.addPath(path)
}

func (this *Point) distanceTo(other Point) float64 {
	return math.Hypot(other.x-this.x, other.y-this.y)
}

func (path *Path) pivotAround(pivot Point, angle float64) {
	sin, cos := math.Sincos(angle)
	for i := range *path {
		(*path)[i].pos.pivotAround(pivot, sin, cos)
	}
}

func (p *Point) pivotAround(pivot Point, sin, cos float64) {
	dx := p.x - pivot.x
	dy := p.y - pivot.y
	p.x = pivot.x + dx*cos - dy*sin
	p.y = pivot.y + dx*sin + dy*cos
}

func (s *Schematic) HandleEntry(entry *parsing.Entry) {
	p1 := Point{s.pos.x, s.pos.y}

	switch entry.Command.Type {
	case "push":
		stack.Push(p1)
		return
	case "pop":
		s.pos = stack.Pop()
		return
	}

	elem := entry.Element

	for _, action := range elem.Actions {

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

	p2 := Point{s.pos.x, s.pos.y}

	switch elem.Type {
	case "resistor":
		s.resistor(p1, p2)
	case "battery":
		s.battery(p1, p2)
	case "line":
		s.line(p1, p2)
	case "capacitor":
		s.capacitor(p1, p2)
	case "dot":
		s.dot(p1)
	default:
		fmt.Printf("unimplemented element type: %s\n", elem.Type)
	}
}

func NewSchematic() *Schematic {
	return &Schematic{}
}

func (s *Schematic) Translate(dx, dy float64) *Schematic {
	s.pos.x += dx
	s.pos.y += dy
	return s
}

// FIXME: this is one ugly function.. 
// we probably need to change the way we're handling svg 
// elements. Can't have duplicated logic for path and circle and god
// knows what more elements...
func (s *Schematic) Normalise() (width, height float64) {
	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64

	// find bounds of paths
	for _, path := range s.paths {
		for _, command := range *path {
			minX = min(minX, command.pos.x)
			minY = min(minY, command.pos.y)
			maxX = max(maxX, command.pos.x)
			maxY = max(maxY, command.pos.y)
		}
	}
	
	// find bounds for circles
	for _, circle := range s.circles {
		minX = min(minX, circle.centerPos.x-circle.radius)
		minY = min(minY, circle.centerPos.y-circle.radius)
		maxX = max(maxX, circle.centerPos.x+circle.radius)
		maxY = max(maxY, circle.centerPos.y+circle.radius)
	}

	// apply translation to all paths
	for _, path := range s.paths {
		for i := range *path {
			(*path)[i].pos.x -= minX - defaultStrokeWidth
			(*path)[i].pos.y -= minY - defaultStrokeWidth
		}
	}

	// apply translation to all circles
	for _, circle := range s.circles {
		circle.centerPos.x -= minX - defaultStrokeWidth
		circle.centerPos.y -= minY - defaultStrokeWidth
	}

	width = maxX - minX + defaultStrokeWidth*2
	height = maxY - minY + defaultStrokeWidth*2
	return width, height
}

func (s *Schematic) End(buf *bytes.Buffer) {
	width, height := s.Normalise()

	buf.WriteString(fmt.Sprintf(
		`<svg width='%d' height='%d' xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1">`,
		int(width), int(height),
	))

	for _, path := range s.paths {
		buf.WriteString(`<path d="`)
		for _, command := range *path {
			buf.WriteString(fmt.Sprintf("%c %g %g ",
				command.letter,
				command.pos.x,
				command.pos.y,
			))
		}
		buf.WriteString(fmt.Sprintf(
			`" style="stroke:black; stroke-width:%d; stroke-linecap: square; fill:none;"></path>`,
			defaultStrokeWidth,
		))
	}

	for _, circle := range s.circles {
		buf.WriteString(fmt.Sprintf(
			`<circle cx="%g" cy="%g" r="%g" fill="white" stroke="black" stroke-width="%d"></circle>`,
			circle.centerPos.x, circle.centerPos.y, circle.radius, defaultStrokeWidth,
		))
	}

	buf.WriteString("</svg>")
}
