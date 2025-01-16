package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sermuns/schemgo/drawing"
	"github.com/sermuns/schemgo/parsing"
)

func exitWithFlagError(message string) {
	fmt.Println("Error:", message)
	flag.Usage()
	os.Exit(1)
}

func main() {
	start := time.Now()

	outputFilePath := flag.String("o", "", "path to output file")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		exitWithFlagError("no input file")
	}
	if len(args)-2 > 1 {
		exitWithFlagError("too many arguments")
	}

	inputFilePath := args[0]

	parsedSchematic := parsing.MustReadSchematic(inputFilePath)
	svgSchematic := drawing.NewSchematic()
	for _, comp := range parsedSchematic.Elements {
		svgSchematic.AddElement(comp)
	}
	svgSchematic.End(*outputFilePath)

	fmt.Printf(
		"Parsed `%s` in %s\n",
		inputFilePath,
		time.Since(start),
	)
}
