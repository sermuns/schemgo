package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sermuns/schemgo/drawing"
	"github.com/sermuns/schemgo/parsing"
)

func writeSchematic(inContents []byte) (outContent []byte) {
	parsedSchematic := parsing.MustReadSchematic(inContents, "")
	svgSchematic := drawing.NewSchematic()
	if len(parsedSchematic.Elements) == 0 {
		fmt.Printf("No elements found in schematic\n")
		os.Exit(1)
	}
	for _, comp := range parsedSchematic.Elements {
		svgSchematic.AddElement(comp)
	}
	var buf bytes.Buffer
	svgSchematic.End(&buf)
	return buf.Bytes()
}
