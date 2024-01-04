/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	R "app.ddcli.datnn/application"
	h "app.ddcli.datnn/application/lib"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
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
		var defaultPort int
		if serverConfig, err := R.AppConfiguration.GetConfig("callbackServer"); err == nil {
			if port, err := serverConfig.Get("port"); err == nil {
				var intPort, _ = strconv.Atoi(port)
				defaultPort = intPort
			}
		}
		var credential = h.Must[*oauth2.Token](h.StartCredentialCallbackServer(defaultPort, &R.AppConfiguration))
		var userInfo = h.HttpGetUserInfo(credential.AccessToken)
		defer userInfo.Body.Close()
		var userInfoResponse h.UserInfoResponse
		if err := json.Unmarshal(h.Must[[]byte](io.ReadAll(userInfo.Body)), &userInfoResponse); err != nil {
			h.LogErr.Printf("Failed to get userinfo %s\n", err.Error())
			os.Exit(1)
		}

		R.GoogleDriveCredential.SetConfig(h.Spr("google%s", userInfoResponse.Id), "refreshToken", credential.RefreshToken)
		R.GoogleDriveCredential.SaveConfig()
		h.LogResult.Printf("Successfully add credential for user %s to application\n", userInfoResponse.Id)
	},
}

var listCredentialCmd = &cobra.Command{
	Use:   "list",
	Short: "Add new Google Drive credential",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Shit happend")
	},
}

func init() {
	credentialCmd.AddCommand(listCredentialCmd)
	credentialCmd.AddCommand(addCredentialCmd)
	rootCmd.AddCommand(credentialCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
