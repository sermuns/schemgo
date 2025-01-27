package cmd

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/sermuns/schemgo/drawing"
	"github.com/sermuns/schemgo/parsing"
)

func performMove(s *drawing.Schematic, move parsing.Movement) {
	value := move.Units * drawing.UnitLength
	if value == 0 {
		value = drawing.DefaultLength
	}

	moveFunc, ok := drawing.MovementFuncs[move.Type]
	if !ok {
		fmt.Printf("unimplemented movement type: %s\n", move.Type)
		os.Exit(1)
	}

	moveFunc(s, value)
}

func renderElem(s *drawing.Schematic, elem *parsing.Element) {
	p1 := s.Pos
	for _, move := range elem.Movements {
		performMove(s, move)
	}
	p2 := s.Pos

	renderFunc, ok := drawing.ElemTypeToRenderFunc[elem.Type]
	if !ok {
		fmt.Printf("unimplemented element type: %s\n", elem.Type)
		os.Exit(1)
	}
	renderFunc(s, p1, p2)

}

func performCommand(s *drawing.Schematic, command parsing.Command) bool {
	if len(command.Type) == 0 {
		return false
	}

	commandFunc, ok := drawing.CommandTypeToFunc[command.Type]
	if !ok {
		fmt.Printf("unimplemented command type: %s\n", command.Type)
		os.Exit(1)
	}
	commandFunc(s)
	return true
}

func writeSchematic(inContents []byte) (outContent []byte) {
	parsedSchematic := parsing.MustReadSchematic(inContents, "")
	s := drawing.NewSchematic()
	if len(parsedSchematic.Entries) == 0 {
		fmt.Printf("No entries in schematic, can't build\n")
		os.Exit(1)
	}

	for _, entry := range parsedSchematic.Entries {
		if performCommand(s, entry.Command) {
			continue
		}
		renderElem(s, &entry.Element)
	}
	var buf bytes.Buffer
	time.Sleep(time.Second)
	s.End(&buf)
	return buf.Bytes()
}
