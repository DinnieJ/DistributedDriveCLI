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
	Use:   "ddc",
	Short: "Distributed Google Drive Storage CLI",
	Long:  `CLI for download and upload file and folder protected to distributed Google Drive Storage`,
}

const CONFIG_FILE_PATH = ".config/ddc/conf.ini"

var AppConfiguration h.Configuration

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var homeDir, _ = os.UserHomeDir()

	var conf = filepath.Join(homeDir, CONFIG_FILE_PATH)
	AppConfiguration.SetFilePath(conf)
	if err := AppConfiguration.Init(); err != nil {
		log.Fatal(err)
	}
	if err := AppConfiguration.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app.ddcli.datnn.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
