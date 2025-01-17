package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [input file]",
	Short: "Build circuit diagram from .schemgo file",
	Args:  cobra.ExactArgs(1),
	Long: `Example:
$ schemgo build examples/simple.schemgo -o simple.svg`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		inFilePath := args[0]
		outFilePath, _ := cmd.Flags().GetString("output")

		inContent, err := os.ReadFile(inFilePath)
		if err != nil {
			fmt.Printf("Error reading file `%s`: %s\n", inFilePath, err)
			os.Exit(1)
		}

		written := writeSchematic(inContent)
		os.WriteFile(outFilePath, written, os.ModePerm)

		fmt.Printf("Parsed `%s` in %s\n", inFilePath, time.Since(start))
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("output", "o", "", "Output file path")
	buildCmd.MarkFlagRequired("output")
}
