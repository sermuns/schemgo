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
	if len(parsedSchematic.Entries) == 0 {
		fmt.Printf("No entries in schematic, can't build\n")
		os.Exit(1)
	}
	for _, entry := range parsedSchematic.Entries{
		svgSchematic.HandleEntry(entry)
	}
	var buf bytes.Buffer
	svgSchematic.End(&buf)
	return buf.Bytes()
}
