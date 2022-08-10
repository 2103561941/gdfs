package cmd

import "github.com/spf13/cobra"

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "rename file in remote file system",
	Args: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
