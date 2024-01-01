/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"net/http"
	"strconv"

	h "app.ddcli.datnn/helpers"
	"github.com/spf13/cobra"
)

// addConfigCmd represents the addConfig command
var addConfigCmd = &cobra.Command{
	Use:   "addConfig",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var defaultPort int
		if serverConfig, err := AppConfiguration.GetConfig("callbackServer"); err == nil {
			if port, err := serverConfig.Get("port"); err == nil {
				var intPort, _ = strconv.Atoi(port)
				defaultPort = intPort
			}
		}

		if err := h.StartCallbackServer(defaultPort); err != nil && errors.Is(err, http.ErrServerClosed) {
			h.LogResult.Println("Add application successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(addConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
