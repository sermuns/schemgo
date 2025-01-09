package main

import (
	"flag"
	"fmt"
	"time"
	"os"

	"github.com/sermuns/schemgo/drawing"
	"github.com/sermuns/schemgo/parsing"
)

const (
	width      = 500
	height     = 500
	outputFile = "index.html"
)

func main() {
	start := time.Now()

	schematicFilePath := flag.String("input", "", "path to .schemgo file")
	flag.Parse()

	// Check if the input flag is provided
	if *schematicFilePath == "" {
		fmt.Println("Error: -input flag is required")
		flag.Usage() // Prints the usage information
		os.Exit(1)   // Exits the program with a non-zero status
	}

	schematic, err := parsing.ReadSchematic(*schematicFilePath)
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
	fmt.Printf("Parsed `%s` in %s\n", *schematicFilePath, time.Since(start))
}
