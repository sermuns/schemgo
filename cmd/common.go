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
	s := drawing.NewSchematic()
	if len(parsedSchematic.Entries) == 0 {
		fmt.Printf("No entries in schematic, can't build\n")
		os.Exit(1)
	}
	for _, entry := range parsedSchematic.Entries {
		p1 := s.Pos

		switch entry.Command.Type {
		case "push":
			s.Push(p1)
			return
		case "pop":
			s.Pos = s.Pop()
			return
		}

		elem := entry.Element

		for _, action := range elem.Actions {

			value := action.Units * drawing.UnitLength
			if value == 0 {
				value = drawing.DefaultLength
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

		p2 := s.Pos

		renderFunc, ok := drawing.ElemTypeToRenderFunc[elem.Type]
		if !ok {
			fmt.Printf("unimplemented element type: %s\n", elem.Type)
			os.Exit(1)
		}
		renderFunc(s, p1, p2)
	}
	var buf bytes.Buffer
	s.End(&buf)
	return buf.Bytes()
}
