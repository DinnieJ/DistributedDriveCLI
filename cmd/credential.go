/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strconv"

	R "app.ddcli.datnn/global"
	h "app.ddcli.datnn/helpers"
	"github.com/spf13/cobra"
)

// addConfigCmd represents the addConfig command
var addConfigCmd = &cobra.Command{
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
		var defaultPort int
		if serverConfig, err := R.AppConfiguration.GetConfig("callbackServer"); err == nil {
			if port, err := serverConfig.Get("port"); err == nil {
				var intPort, _ = strconv.Atoi(port)
				defaultPort = intPort
			}
		}

		if cred, err := h.StartCredentialCallbackServer(defaultPort, &R.AppConfiguration); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(cred)
		}
	},
}

func init() {
	addConfigCmd.AddCommand(addCredentialCmd)
	rootCmd.AddCommand(addConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
