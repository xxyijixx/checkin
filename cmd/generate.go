package cmd

import (
	generate "checkin/query/gen"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate code",
	Run: func(cmd *cobra.Command, args []string) {
		generate.Generate()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
