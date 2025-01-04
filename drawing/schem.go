package drawing

import (
	"bytes"
	"os"

	"github.com/ajstarks/svgo"
)

type Schematic struct {
	X, Y   int // Current position of the brush
	Color  string
	Canvas *svg.SVG
	Buffer *bytes.Buffer
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

func (s *Schematic) ChangeBrushColor(color string) *Schematic {
	s.Color = color
	return s
}

func (s *Schematic) Circle(radius int) *Schematic {
	s.Canvas.Circle(s.X, s.Y, radius, "fill:"+s.Color)
	return s
}

func (s *Schematic) Left(units int) *Schematic {
	s.X -= units
	return s
}

func (s *Schematic) Right(units int) *Schematic {
	s.X += units
	return s
}

// Up moves the brush up and updates the position, supports chaining.
func (s *Schematic) Up(units int) *Schematic {
	s.Y -= units
	return s
}

func (s *Schematic) Down(units int) *Schematic {
	s.Y += units
	return s
}

func (s *Schematic) End(outputFile string) error {
	s.Canvas.End()
	return os.WriteFile(outputFile, s.Buffer.Bytes(), 0644)
}
