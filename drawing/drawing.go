package drawing

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"

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

type Path []command

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
			newPath = append(newPath, newCommand)
		}
	}
	return newPath
}

func (s *Schematic) addPath(path Path) {
	s.paths = append(s.paths, &path)
}

func (s *Schematic) addAndPivotPath(start, end point, path Path) {
	angle := math.Atan2(end.y-start.y, end.x-start.x)
	path.pivotAround(start, angle)
	s.addPath(path)
}

func (this *point) distanceTo(other point) float64 {
	return math.Hypot(other.x-this.x, other.y-this.y)
}

func (path *Path) pivotAround(pivot point, angle float64) {
	sin, cos := math.Sincos(angle)
	for i := range *path {
		(*path)[i].pos.pivotAround(pivot, sin, cos)
	}
}

func (p *point) pivotAround(pivot point, sin, cos float64) {
	dx := p.x - pivot.x
	dy := p.y - pivot.y
	p.x = pivot.x + dx*cos - dy*sin
	p.y = pivot.y + dx*sin + dy*cos
}

func (s *Schematic) resistor(p1, p2 point) {
	const (
		height = defaultLength / 4
		width  = defaultLength / 2
	)

	distance := p1.distanceTo(p2)

	s.addAndPivotPath(p1, p2, createPath(
		false, p1.x, p1.y,
		true, p1.x+distance/2-width/2, p1.y,
		true, p1.x+distance/2-width/2, p1.y+height/2,
		true, p1.x+distance/2+width/2, p1.y+height/2,
		true, p1.x+distance/2+width/2, p1.y-height/2,
		true, p1.x+distance/2-width/2, p1.y-height/2,
		true, p1.x+distance/2-width/2, p1.y,
		false, p1.x+distance/2+width/2, p1.y,
		true, p1.x+distance, p1.y,
	))
}

func (s *Schematic) battery(p1, p2 point) {
	const (
		termGap       = defaultLength / 6
		negTermHeight = defaultLength / 4
		posTermHeight = defaultLength / 2
	)

	distance := p1.distanceTo(p2)
	negTermX := p1.x + distance/2 - termGap/2
	posTermX := negTermX + termGap

	s.addAndPivotPath(p1, p2, createPath(
		false, p1.x, p1.y,
		true, negTermX, p1.y,
		false, negTermX, p1.y-negTermHeight/2,
		true, negTermX, p1.y+negTermHeight/2,
		false, posTermX, p1.y-posTermHeight/2,
		true, posTermX, p1.y+posTermHeight/2,
		false, posTermX, p1.y,
		true, p1.x+distance, p1.y,
	))

}

func (s *Schematic) AddElement(e *parsing.Element) {
	p1 := point{s.X, s.Y}

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

	p2 := point{s.X, s.Y}

	switch e.Type {
	case "resistor":
		s.resistor(p1, p2)
	case "battery":
		s.battery(p1, p2)
	case "line":
		s.line(p1, p2)
	}
}

func NewSchematic() *Schematic {
	return &Schematic{
		X:     0,
		Y:     0,
		paths: []*Path{},
	}
}

func (s *Schematic) line(p1, p2 point) {
	s.addPath(createPath(
		false, p1.x, p1.y,
		true, p2.x, p2.y,
	))
}

func (s *Schematic) Translate(dx, dy float64) *Schematic {
	s.X += dx
	s.Y += dy
	return s
}

func (s *Schematic) Normalise() (width, height float64) {
	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64

	// find bounds
	for _, path := range s.paths {
		for _, command := range *path {
			minX = min(minX, command.pos.x)
			minY = min(minY, command.pos.y)
			maxX = max(maxX, command.pos.x)
			maxY = max(maxY, command.pos.y)
		}
	}

	// apply translation to all paths
	for _, path := range s.paths {
		for i := range *path {
			(*path)[i].pos.x -= minX - defaultStrokeWidth
			(*path)[i].pos.y -= minY - defaultStrokeWidth
		}
	}

	width = maxX - minX + defaultStrokeWidth*2
	height = maxY - minY + defaultStrokeWidth*2
	return width, height
}

// need to create root <svg> tag
// and convert all paths to <path d=...> tags
func (s *Schematic) End(outFilePath string) {
	var buf bytes.Buffer

	width, height := s.Normalise()

	buf.WriteString(fmt.Sprintf(
		`<svg width='%d' height='%d' xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1">`,
		int(width), int(height),
	))

	for _, path := range s.paths {
		buf.WriteString(`<path `)
		buf.WriteString(` d="`)
		for _, command := range *path {
			if command.penDown {
				buf.WriteString("L ")
			} else {
				buf.WriteString("M ")
			}
			buf.WriteString(fmt.Sprintf("%d %d ",
				int(command.pos.x),
				int(command.pos.y),
			))
		}
		buf.WriteString(`" style="stroke:black; stroke-width:5; stroke-linecap: square; fill:none;"></path>`)
	}
	buf.WriteString("</svg>")

	os.MkdirAll(filepath.Dir(outFilePath), os.ModePerm)
	os.WriteFile(outFilePath, buf.Bytes(), os.ModePerm)
}
