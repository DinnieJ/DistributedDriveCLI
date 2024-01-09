/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	h "app.ddcli.datnn/application/lib"
	"github.com/spf13/cobra"
)

// addConfigCmd represents the addConfig command
var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "Credential group command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		h.AddZip("application/impl/credential.go", "application.zip")
	},
}

func init() {
	rootCmd.AddCommand(zipCmd)
}
