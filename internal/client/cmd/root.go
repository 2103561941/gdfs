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


// 开启命令行服务
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// 自动生成markdown使用文档
func GenDocs(filepath string) error {
	return doc.GenMarkdownTree(rootCmd, filepath)
}
