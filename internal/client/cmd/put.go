package cmd

import "github.com/spf13/cobra"

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "put local file to remote file system",
	Args: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(putCmd)
}
