package root

import (
	"log"
	"os"
	"path/filepath"

	h "app.ddcli.datnn/application/lib"
)

const CONFIG_FILE_PATH = ".config/ddc/ddcli.conf"
const CREDENTIAL_FILE_PATH = ".config/ddc/credential.conf"

var AppConfiguration h.Configuration
var StoreCredential h.Configuration

func init() {
	var homeDir, _ = os.UserHomeDir()

	var conf = filepath.Join(homeDir, CONFIG_FILE_PATH)
	var credentialPath = filepath.Join(homeDir, CREDENTIAL_FILE_PATH)
	AppConfiguration.SetFilePath(conf)
	StoreCredential.SetFilePath(credentialPath)
	StoreCredential.Generated = true
	if err := h.LoadAllConfig(&AppConfiguration, &StoreCredential); err != nil {
		log.Fatal(err)
	}
}
