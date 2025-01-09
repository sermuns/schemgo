package main

import (
	"github.com/sermuns/schemgo/drawing"
	// "fmt"

	"github.com/sermuns/schemgo/parsing"
)

const (
	width             = 500
	height            = 500
	outputFile        = "index.html"
	schematicFilePath = "./examples/simple.schemgo"
)

func main() {
	schematic, err := parsing.ReadSchematic(schematicFilePath)
	if err != nil {
		panic(err)
	}

	s := drawing.NewSchematic(width, height)

	for _, comp := range schematic.Elements {
		s.AddElement(comp)
	}

	err = s.End(outputFile)
	if err != nil {
		panic(err)
	}
}
