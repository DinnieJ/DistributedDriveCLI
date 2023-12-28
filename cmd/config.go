/*
Copyright © 2023 Datnn <datnn288@gmail.com>
*/
package cmd

import (
	"fmt"

	h "app.ddcli.datnn/helpers"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create config for Google Drive Application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var clientId = h.GetInput("Google Client ID: ")
		var clientSecret = h.GetInput("Google Client Secret: ")
		AppConfiguration.SetConfig("googleAppClient", "clientID", clientId)
		AppConfiguration.SetConfig("googleAppClient", "clientSecret", clientSecret)
		AppConfiguration.WriteToConfigFile()
		color.Blue("Setup config successful")
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all config for application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(AppConfiguration.GetPrtString())
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	rootCmd.AddCommand(configCmd)
}
