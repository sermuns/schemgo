package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/sermuns/schemgo/drawing"
	"github.com/sermuns/schemgo/parsing"
	"github.com/spf13/cobra"
)

func writeSchematic(inputFileContents []byte, inputFilePath, outputFilePath string) {
	parsedSchematic := parsing.MustReadSchematic(inputFileContents, inputFilePath)
	svgSchematic := drawing.NewSchematic()
	if(len(parsedSchematic.Elements) == 0) {
		fmt.Printf("No elements found in `%s`\n", inputFilePath)
		os.Exit(1)
	}
	for _, comp := range parsedSchematic.Elements {
		svgSchematic.AddElement(comp)
	}
	svgSchematic.End(outputFilePath)
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [input file]",
	Short: "Build circuit diagram from .schemgo file",
	Args:  cobra.ExactArgs(1),
	Long: `Example:
$ schemgo build examples/simple.schemgo -o simple.svg`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		inputFilePath := args[0]

		if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
			fmt.Printf("File `%s` does not exist\n", inputFilePath)
			os.Exit(1)
		}

		inputFileContents, err := os.ReadFile(inputFilePath)
		if err != nil {
			fmt.Printf("Error reading file `%s`\n", inputFilePath)
			os.Exit(1)
		}

		outputFilePath, err := cmd.Flags().GetString("output")

		writeSchematic(inputFileContents, inputFilePath, outputFilePath)

		fmt.Printf("Parsed `%s` in %s\n", inputFilePath, time.Since(start))
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.ValidArgs = []string{"input"}

	buildCmd.Flags().StringP("output", "o", "", "Output file path")
	buildCmd.MarkFlagRequired("output")
}
