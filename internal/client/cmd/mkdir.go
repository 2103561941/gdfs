package cmd

import "github.com/spf13/cobra"

var mkdirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "make directory in remote file system",
	Args: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(mkdirCmd)
}
