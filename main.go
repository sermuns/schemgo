package main

import (
	"bytes"
	"log"
	"os"

	"github.com/ajstarks/svgo"
)

const (
	width      = 500
	height     = 500
	outputFile = "out.svg"
)

func main() {
	var outputBuffer bytes.Buffer

	canvas := svg.New(&outputBuffer)
	canvas.Start(width, height)
	canvas.Circle(width/2, height/2, 100)
	style := `
		fill: white;
		dominant-baseline: central;
		text-anchor: middle;
		font-size: 30px;
		font-family: New Computer Modern Math;
	`
	canvas.Text(width/2, height/2, "Hello, SVG", style)
	canvas.End()

	file, _ := os.Create(outputFile)

	file.Write(outputBuffer.Bytes())
	log.Printf("Written %d bytes\n", outputBuffer.Len())
}
