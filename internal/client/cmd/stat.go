package cmd

import "github.com/spf13/cobra"

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "get metadata of remote file system",
	Args: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

func init() {
	rootCmd.AddCommand(statCmd)
}
