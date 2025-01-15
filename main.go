package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sermuns/schemgo/drawing"
	"github.com/sermuns/schemgo/parsing"
)

func main() {
	start := time.Now()

	schematicFilePath := flag.String("input", "", "path to .schemgo")
	flag.StringVar(schematicFilePath, "i", "", "shorthand for input")

	outputFilePath := flag.String("output", "", "path to output")
	flag.StringVar(outputFilePath, "o", "", "shorthand for output")

	flag.Parse()

	if *schematicFilePath == "" {
		fmt.Println("Error: input flag is required")
		flag.Usage() // Prints the usage information
		os.Exit(1)   // Exits the program with a non-zero status
	}

	schematic, err := parsing.ReadSchematic(*schematicFilePath)
	if err != nil {
		panic(err)
	}

	s := drawing.NewSchematic()

	for _, comp := range schematic.Elements {
		// fmt.Printf("Adding %s\n", comp.Type)
		s.AddElement(comp)
	}

	s.End(*outputFilePath)

	fmt.Printf("Parsed `%s` in %s\n", *schematicFilePath, time.Since(start))
}
