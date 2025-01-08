package main

import (
	"github.com/sermuns/schemgo/drawing"
	// "github.com/sermuns/schemgo/parsing"
)

const (
	width             = 500
	height            = 500
	outputFile        = "index.html"
	schematicFilePath = "./examples/simple.schemgo"
)

func main() {
	s := drawing.NewSchematic(width, height).
		ChangeBrushColor("red").
		Circle(50).Left(100).
		ChangeBrushColor("blue").
		Circle(30).Right(200)

	s.ChangeCanvasSize(600, 1000)

	err := s.End(outputFile)

	if err != nil {
		panic(err)
	}

	// schematic := parsing.ReadSchematic(schematicFilePath)
}
