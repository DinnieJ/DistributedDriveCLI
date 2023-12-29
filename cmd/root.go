/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	h "app.ddcli.datnn/helpers"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ddcli",
	Short: "Distributed Google Drive Storage CLI",
	Long:  `CLI for download and upload file and folder protected to distributed Google Drive Storage`,
}

const CONFIG_FILE_PATH = ".config/ddc/ddcli.conf"
const CREDENTIAL_FILE_PATH = ".config/ddc/credential.conf"

var AppConfiguration h.Configuration
var GoogleDriveCredential h.Configuration

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var homeDir, _ = os.UserHomeDir()

	var conf = filepath.Join(homeDir, CONFIG_FILE_PATH)
	var credentialPath = filepath.Join(homeDir, CREDENTIAL_FILE_PATH)
	AppConfiguration.SetFilePath(conf)
	GoogleDriveCredential.SetFilePath(credentialPath)
	if err := h.LoadAllConfig(&AppConfiguration, &GoogleDriveCredential); err != nil {
		log.Fatal(err)
	}
}
