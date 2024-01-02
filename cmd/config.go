/*
Copyright © 2023 Datnn <datnn288@gmail.com>
*/
package cmd

import (
	"fmt"

	h "app.ddcli.datnn/helpers"
	R "app.ddcli.datnn/root"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create config for Google Drive Application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var clientId = h.GetInput("Google Client ID: ")
		var clientSecret = h.GetInput("Google Client Secret: ")
		R.AppConfiguration.SetConfig("googleAppClient", "clientID", clientId)
		R.AppConfiguration.SetConfig("googleAppClient", "clientSecret", clientSecret)
		R.AppConfiguration.SaveConfig()
		h.LogResult.Println("Setup config successful")
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all config for application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(R.AppConfiguration.GetPrtString())
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	rootCmd.AddCommand(configCmd)
}
