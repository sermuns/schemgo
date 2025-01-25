package cmd

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Root   string `json:"root"`
	Config struct {
		Book struct {
			Authors      []string `json:"authors"`
			Language     string   `json:"language"`
			Multilingual bool     `json:"multilingual"`
			Src          string   `json:"src"`
			Title        string   `json:"title"`
		} `json:"book"`
		Preprocessor map[string]struct {
			Command string `json:"command"`
		} `json:"preprocessor"`
	} `json:"config"`
	Renderer      string `json:"renderer"`
	MdbookVersion string `json:"mdbook_version"`
}

type Section struct {
	Chapter struct {
		Name        string   `json:"name"`
		Content     string   `json:"content"`
		Number      []int    `json:"number"`
		SubItems    []string `json:"sub_items"`
		Path        string   `json:"path"`
		SourcePath  string   `json:"source_path"`
		ParentNames []string `json:"parent_names"`
	} `json:"Chapter"`
}

type Data struct {
	Sections      []Section `json:"sections"`
	NonExhaustive any       `json:"__non_exhaustive"`
}

var mdbookCmd = &cobra.Command{
	Use:   "mdbook",
	Short: "Act as mdBook preprocessor. You probably don't want to manually use this!",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 && args[1] == "supports" {
			os.Exit(0)
		}

		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read stdin: %v\n", err)
			os.Exit(1)
		}

		os.WriteFile("input.json", stdin, 0644)
		return

		var input []any
		err = json.Unmarshal(stdin, &input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to unmarshal stdin: %v\n", err)
			os.Exit(1)
		}

		book, err := json.Marshal(input[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to marshal book: %v\n", err)
			os.Exit(1)
		}

		var data Data
		err = json.Unmarshal(book, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to unmarshal book: %v\n", err)
			os.Exit(1)
		}

		for i := range data.Sections {
			data.Sections[i].Chapter.Content = "Hello, World!"
		}

		output, err := json.Marshal(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to marshal updated data: %v\n", err)
			os.Exit(1)
		}

		os.WriteFile("output.json", output, 0644)

	},
}

func init() {
	rootCmd.AddCommand(mdbookCmd)
}
