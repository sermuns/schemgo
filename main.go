package main

import (
	"github.com/sermuns/schemgo/drawing"
)

const (
	width      = 500
	height     = 500
	outputFile = "index.html"
)

func main() {
	s := drawing.NewSchematic(width, height).
		ChangeBrushColor("red").
		Circle(50).Left(100).
		ChangeBrushColor("blue").
		Circle(30).Right(200)

	err := s.End(outputFile)
	if err != nil {
		panic(err)
	}
}
