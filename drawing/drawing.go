package drawing

import (
	"bytes"
	"fmt"
	"math"
)

const (
	UnitLength         = 100
	DefaultLength      = 1 * UnitLength
	DefaultStrokeWidth = 5
)

type Point struct {
	X, Y float64
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
	posStack Stack
	Pos      Point
	Paths    []*Path
	Circles  []*Circle
}

func (s *Schematic) Push(p Point) {
	s.posStack.Push(s.Pos)
}

func (s *Schematic) Pop() Point {
	return s.posStack.Pop()
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
			newCommand.pos.X = part.(float64)
		case 2:
			newCommand.pos.Y = part.(float64)
			newPath = append(newPath, newCommand)
		}
	}
	return newPath
}

func (s *Schematic) addPath(path Path) {
	s.Paths = append(s.Paths, &path)
}

func (s *Schematic) addCircle(x, y, radius float64) {
	s.Circles = append(s.Circles, &Circle{
		centerPos: Point{x, y},
		radius:    radius,
	})
}

func (s *Schematic) addAndPivotPath(start, end Point, path Path) {
	angle := math.Atan2(end.Y-start.Y, end.X-start.X)
	path.pivotAround(start, angle)
	s.addPath(path)
}

func (this *Point) distanceTo(other Point) float64 {
	return math.Hypot(other.X-this.X, other.Y-this.Y)
}

func (path *Path) pivotAround(pivot Point, angle float64) {
	sin, cos := math.Sincos(angle)
	for i := range *path {
		(*path)[i].pos.pivotAround(pivot, sin, cos)
	}
}

func (p *Point) pivotAround(pivot Point, sin, cos float64) {
	dx := p.X - pivot.X
	dy := p.Y - pivot.Y
	p.X = pivot.X + dx*cos - dy*sin
	p.Y = pivot.Y + dx*sin + dy*cos
}

func NewSchematic() *Schematic {
	return &Schematic{}
}

func (s *Schematic) Translate(dx, dy float64) *Schematic {
	s.Pos.X += dx
	s.Pos.Y += dy
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
	for _, path := range s.Paths {
		for _, command := range *path {
			minX = min(minX, command.pos.X)
			minY = min(minY, command.pos.Y)
			maxX = max(maxX, command.pos.X)
			maxY = max(maxY, command.pos.Y)
		}
	}

	// find bounds for circles
	for _, circle := range s.Circles {
		minX = min(minX, circle.centerPos.X-circle.radius)
		minY = min(minY, circle.centerPos.Y-circle.radius)
		maxX = max(maxX, circle.centerPos.X+circle.radius)
		maxY = max(maxY, circle.centerPos.Y+circle.radius)
	}

	// apply translation to all paths
	for _, path := range s.Paths {
		for i := range *path {
			(*path)[i].pos.X -= minX - DefaultStrokeWidth
			(*path)[i].pos.Y -= minY - DefaultStrokeWidth
		}
	}

	// apply translation to all circles
	for _, circle := range s.Circles {
		circle.centerPos.X -= minX - DefaultStrokeWidth
		circle.centerPos.Y -= minY - DefaultStrokeWidth
	}

	width = maxX - minX + DefaultStrokeWidth*2
	height = maxY - minY + DefaultStrokeWidth*2
	return width, height
}

func (s *Schematic) End(buf *bytes.Buffer) {
	width, height := s.Normalise()

	buf.WriteString(fmt.Sprintf(
		`<svg width='%d' height='%d' xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1">`,
		int(width), int(height),
	))

	for _, path := range s.Paths {
		buf.WriteString(`<path d="`)
		for _, command := range *path {
			buf.WriteString(fmt.Sprintf("%c %g %g ",
				command.letter,
				command.pos.X,
				command.pos.Y,
			))
		}
		buf.WriteString(fmt.Sprintf(
			`" style="stroke:black; stroke-width:%d; stroke-linecap: square; fill:none;"></path>`,
			DefaultStrokeWidth,
		))
	}

	for _, circle := range s.Circles {
		buf.WriteString(fmt.Sprintf(
			`<circle cx="%g" cy="%g" r="%g" fill="white" stroke="black" stroke-width="%d"></circle>`,
			circle.centerPos.X, circle.centerPos.Y, circle.radius, DefaultStrokeWidth,
		))
	}

	buf.WriteString("</svg>")
}
