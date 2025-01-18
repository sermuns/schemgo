package cmd

import (
	// "encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var mdbookCmd = &cobra.Command{
	Use:   "mdbook",
	Short: "Act as mdBook preprocessor. You probably don't want to manually use this!",
	Run: func(cmd *cobra.Command, args []string) {
		// mdbook wants this
		if len(args) > 1 && args[1] == "supports" {
			os.Exit(0)
		}

		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			os.Exit(1)
		}

		// var data map[string]interface{}

		if err != nil {
			fmt.Println("heh", content)
			fmt.Println("Error parsing JSON:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(mdbookCmd)
}
