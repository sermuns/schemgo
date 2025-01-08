package drawing

import (
	"bytes"
	"fmt"
	"os"
	"strings"

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

func (s *Schematic) Left(units int) *Schematic {
	s.X -= units
	return s
}

func (s *Schematic) Right(units int) *Schematic {
	s.X += units
	return s
}

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
