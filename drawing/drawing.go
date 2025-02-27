package drawing

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
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
	points []Point
}

type Path []command

type Circle struct {
	centerPos Point
	radius    float64
	style     string
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

// i'm sorry this is so fucked up
func createPath(input ...any) Path {
	var path Path
	var thisCommand command

	for _, part := range input {
		switch v := part.(type) {
		case rune: // If it's a command letter
			// Append the current command to the path if it has a letter
			if thisCommand.letter != 0 {
				path = append(path, thisCommand)
			}
			// Start a new command
			thisCommand = command{letter: v}

		case Point: // If it's a point
			if thisCommand.letter == 0 {
				fmt.Println("createPath: point provided before command letter")
				os.Exit(1)
			}
			thisCommand.points = append(thisCommand.points, v)

		default: // Unknown type
			fmt.Printf("createPath: unknown type %T\n", part)
			os.Exit(1)
		}
	}

	// Append the final command if valid
	if thisCommand.letter != 0 {
		path = append(path, thisCommand)
	}

	return path
}

func (s *Schematic) addPath(path Path) {
	s.Paths = append(s.Paths, &path)
}

func (s *Schematic) addCircle(x, y, radius float64, optionalStyle ...string) {
	var style string
	if len(optionalStyle) > 0 {
		style = optionalStyle[0]
	} else {
		style = fmt.Sprintf(`fill="black" stroke="black" stroke-width="%d"`, DefaultStrokeWidth)
	}

	s.Circles = append(s.Circles, &Circle{
		centerPos: Point{x, y},
		radius:    radius,
		style:     style,
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

func (path *Path) pivotAround(pivotCenter Point, angle float64) {
	sin, cos := math.Sincos(angle)
	for i, pathCommand := range *path {
		// lowercase indicates relative
		if unicode.IsLower(pathCommand.letter) {
			for j := range (*path)[i].points {
				(*path)[i].points[j].pivotAround(Point{0, 0}, sin, cos)
			}
		} else {
			for j := range (*path)[i].points {
				(*path)[i].points[j].pivotAround(pivotCenter, sin, cos)
			}
		}
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
func (s *Schematic) Normalise() (width, height float64) {

	const padding = 5

	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64

	// find bounds of paths
	for _, path := range s.Paths {
		for _, command := range *path {
			if unicode.IsLower(command.letter) {
				startPoint := command.points[0]
				for _, point := range command.points[1:] {
					startPoint.X += point.X
					startPoint.Y += point.Y
					minX = min(minX, startPoint.X)
					minY = min(minY, startPoint.Y)
					maxX = max(maxX, startPoint.X)
					maxY = max(maxY, startPoint.Y)
				}
			} else {
				for _, point := range command.points {
					minX = min(minX, point.X)
					minY = min(minY, point.Y)
					maxX = max(maxX, point.X)
					maxY = max(maxY, point.Y)
				}
			}
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
		for i, command := range *path {
			if unicode.IsLower(command.letter) {
				continue
			}
			for j := range (*path)[i].points {
				(*path)[i].points[j].X -= minX - DefaultStrokeWidth - padding
				(*path)[i].points[j].Y -= minY - DefaultStrokeWidth - padding
			}
		}
	}

	// apply translation to all circles
	for _, circle := range s.Circles {
		circle.centerPos.X -= minX - DefaultStrokeWidth - padding
		circle.centerPos.Y -= minY - DefaultStrokeWidth - padding
	}

	width = maxX - minX + DefaultStrokeWidth*2 + padding*2
	height = maxY - minY + DefaultStrokeWidth*2 + padding*2
	return width, height
}

func (c *command) asString() string {
	var sb strings.Builder

	sb.WriteByte(byte(c.letter))
	for _, point := range c.points {
		sb.WriteByte(' ')
		sb.Write(strconv.AppendFloat(nil, point.X, 'f', -1, 64))
		sb.WriteByte(',')
		sb.Write(strconv.AppendFloat(nil, point.Y, 'f', -1, 64))
		sb.WriteByte(' ')
	}

	return sb.String()
}

func (s *Schematic) End() []byte {
	width, height := s.Normalise()

	var buf bytes.Buffer

	buf.WriteString(`<svg width="`)
	buf.Write(strconv.AppendInt(nil, int64(width), 10))
	buf.WriteString(`" height="`)
	buf.Write(strconv.AppendInt(nil, int64(height), 10))
	buf.WriteString(`" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1">`)

	// Write paths
	for _, path := range s.Paths {
		buf.WriteString(`<path d="`)
		for _, pathCommand := range *path {
			buf.WriteString(pathCommand.asString())
		}
		buf.WriteString(`" style="stroke:black; stroke-width:`)
		buf.Write(strconv.AppendInt(nil, int64(DefaultStrokeWidth), 10))
		buf.WriteString(`; stroke-linecap: square; fill:none;"></path>`)
	}

	// Write circles
	for _, circle := range s.Circles {
		buf.WriteString(`<circle cx="`)
		buf.Write(strconv.AppendFloat(nil, circle.centerPos.X, 'f', -1, 64))
		buf.WriteString(`" cy="`)
		buf.Write(strconv.AppendFloat(nil, circle.centerPos.Y, 'f', -1, 64))
		buf.WriteString(`" r="`)
		buf.Write(strconv.AppendFloat(nil, circle.radius, 'f', -1, 64))
		buf.WriteString(`" `)
		buf.WriteString(circle.style)
		buf.WriteString(`></circle>`)
	}

	buf.WriteString("</svg>")

	return buf.Bytes()
}
