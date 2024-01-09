/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	f "app.ddcli.datnn/application/impl"
	"github.com/spf13/cobra"
)

// addConfigCmd represents the addConfig command
var credentialCmd = &cobra.Command{
	Use:   "credential",
	Short: "Credential group command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

var addCredentialCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new Google Drive credential",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		f.AddNewStorage()
	},
}

var listCredentialCmd = &cobra.Command{
	Use:   "list",
	Short: "Add new Google Drive credential",
	Run: func(cmd *cobra.Command, args []string) {
		f.GetAllStorageInformation()
	},
}

func init() {
	credentialCmd.AddCommand(listCredentialCmd)
	credentialCmd.AddCommand(addCredentialCmd)
	rootCmd.AddCommand(credentialCmd)
}
