package cmd

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list files in remote file system",
	Args: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
