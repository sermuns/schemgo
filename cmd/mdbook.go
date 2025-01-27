package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type Chapter struct {
	Name        string   `json:"name"`
	Content     string   `json:"content"`
	Number      []int    `json:"number"`
	SubItems    []string `json:"sub_items"`
	Path        string   `json:"path"`
	SourcePath  string   `json:"source_path"`
	ParentNames []string `json:"parent_names"`
}

type Section struct {
	Chapter Chapter `json:"Chapter"`
}

type Config struct {
	Sections      []Section `json:"sections"`
	NonExhaustive *string   `json:"__non_exhaustive"`
}

// find ```schemgo ``` md code blocks, run them through the schemgo processor.
// replace entire block (including ticks) with the output of processor
func processSection(section *Section){
	// go throuhg line-by-line, find
	// ```schemgo
	

}

var mdbookCmd = &cobra.Command{
	Use:   "mdbook",
	Short: "Act as mdBook preprocessor. You probably don't want to manually use this!",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 && args[0] == "supports" {
			os.Exit(0)
		}

		stdInput, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read stdin: %v\n", err)
			os.Exit(1)
		}

		var jsonData []Config
		err = json.Unmarshal(stdInput, &jsonData)
		if err != nil {
			log.Fatal(err)
		}

		book := jsonData[1]
		for i := range book.Sections{
			processSection(&book.Sections[i])
		}

		modifiedJson, err := json.MarshalIndent(jsonData[1], "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(modifiedJson))
	},
}

func init() {
	rootCmd.AddCommand(mdbookCmd)
}
