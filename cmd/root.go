package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func handlePiped() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("Error checking stdin:", err)
		os.Exit(1)
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			os.Exit(1)
		}
		fmt.Print(string(writeSchematic(content)))
		return true
	}
	return false
}

var rootCmd = &cobra.Command{
	Use:   "schemgo",
	Short: "Dead simple circuit schematic generator",
	Long: `Input can also be piped into this root command to get output in stdout.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if handlePiped() {
			return
		}
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
