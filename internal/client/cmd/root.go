/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "gdfs client command line tool",
}


// Register command server  
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func GenDocs(filepath string) error {
	return doc.GenMarkdownTree(rootCmd, filepath)
}
