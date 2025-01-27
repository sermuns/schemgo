package cmd

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkBuild(b *testing.B) {
	inFilePath := "../examples/simple.schemgo"

	inContent, err := os.ReadFile(inFilePath)
	if err != nil {
		fmt.Printf("Error reading file `%s`: %s\n", inFilePath, err)
		os.Exit(1)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		writeSchematic(inContent)
	}
}
