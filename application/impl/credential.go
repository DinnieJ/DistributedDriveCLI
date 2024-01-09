package impl

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	R "app.ddcli.datnn/application"
	h "app.ddcli.datnn/application/lib"
	s "app.ddcli.datnn/application/services"
	sDrive "app.ddcli.datnn/application/services/google-drive"
)

func AddNewStorage() {
	var defaultPort int
	if serverConfig, err := R.AppConfiguration.GetConfig("callbackServer"); err == nil {
		if port, err := serverConfig.Get("port"); err == nil {
			var intPort, _ = strconv.Atoi(port)
			defaultPort = intPort
		}
	}
	var credential = h.Must(h.StartCredentialCallbackServer(defaultPort, &R.AppConfiguration))
	var userInfo = h.HttpGetUserInfo(credential.AccessToken)
	defer userInfo.Body.Close()
	var userInfoResponse h.UserInfoResponse
	if err := json.Unmarshal(h.Must[[]byte](io.ReadAll(userInfo.Body)), &userInfoResponse); err != nil {
		h.LogErr.Printf("Failed to get userinfo %s\n", err.Error())
		os.Exit(1)
	}

	R.StoreCredential.SetConfig(h.Spr("GoogleDrive|%s", userInfoResponse.Id), "refreshToken", credential.RefreshToken)
	R.StoreCredential.SetConfig(h.Spr("GoogleDrive|%s", userInfoResponse.Id), "email", userInfoResponse.Email)
	R.StoreCredential.SaveConfig()
	h.LogResult.Printf("Successfully add credential for user %s to application\n", userInfoResponse.Id)
}

func getStorageService(config *h.Config) s.ServiceResource {
	var serviceResource s.ServiceResource
	var re = regexp.MustCompile(`^(?P<ServiceName>[A-Za-z0-9]+)|.*$`)
	var matches = re.FindStringSubmatch(config.ConfigName)
	if idx := re.SubexpIndex("ServiceName"); idx != -1 {
		switch matches[idx] {
		case "GoogleDrive":
			serviceResource = &sDrive.GoogleDriveService{}
		default:
			serviceResource = nil
		}

	}
	serviceResource.New(config)
	return serviceResource
}

type resultStorageInformation struct {
	data *s.StorageInformation
	err  error
}

// getStorageBar generates a storage bar based on the given value and range.
//
// Parameters:
// - value: a float64 representing the value of the storage.
// - rng: a float64 representing the range of the storage.
//
// Return type: a string representing the storage bar.
func getStorageBar(value float64, rng float64) string {
	if rng == 0 || value > rng {
		return ""
	}
	var storagePercent = value / rng
	var barLength = 30
	var usedBarLength = math.RoundToEven(storagePercent * float64(barLength))
	var bar = ""
	bar += strings.Repeat("█", int(usedBarLength))
	bar += strings.Repeat("░", barLength-int(usedBarLength))
	bar += " "
	bar += fmt.Sprintf("%.2f %%", storagePercent*100)
	return bar
}

// byteToVerboseString converts a byte value into a human-readable string representation.
//
// It takes a float64 value representing the byte value and returns a string.
func byteToVerboseString(value float64) string {
	var unit = []string{
		"KiB", "MiB", "GiB", "TiB",
	}
	var result = h.Spr("%.2f bytes", value)
	for i, v := range unit {
		var n = math.Pow(1024, float64(i+1))
		if value >= n {
			result = h.Spr("%.2f %s", value/n, v)
			continue
		} else {
			break
		}
	}

	return result
}

// GetAllStorageInformation retrieves storage information from multiple sources and logs the results.
//
// If no storage is found, logs a message and exits the program.
// For each storage credential, retrieves the configuration and the corresponding storage service.
// It launches a goroutine for each service to asynchronously retrieve storage information.
// The retrieved information is appended to a list.
// Once all goroutines have completed, it calculates the total storage and usage.
// It logs the storage information and usage bar for each storage.

func GetAllStorageInformation() {
	var wg = &sync.WaitGroup{}
	var mutex = &sync.Mutex{}
	var lStore = len(R.StoreCredential.Data)
	if lStore == 0 {
		h.LogAbout.Println("No storage found")
		os.Exit(0)
	}
	wg.Add(lStore)
	var listResult = make([]resultStorageInformation, 0, lStore)
	for k := range R.StoreCredential.Data {
		var config, err = R.StoreCredential.GetConfig(k)
		if err != nil {
			h.LogErr.Println("Error: Can't find config", k)
			continue
		}
		var service = getStorageService(config)
		go func(s s.ServiceResource) {
			var info, err = s.GetInformation()
			mutex.Lock()
			var result = &resultStorageInformation{
				data: info,
				err:  err,
			}
			listResult = append(listResult, *result)
			mutex.Unlock()
			wg.Done()
		}(service)
	}
	wg.Wait()
	var totalStorage float64 = 0
	var totalUsage float64 = 0
	for _, result := range listResult {
		if result.err != nil {
			h.LogErr.Println("[ERROR]:", result.err.Error())
			continue
		}
		totalStorage += result.data.TotalStorage
		totalUsage += result.data.UsedStorage
		h.LogWarning.Print(getStorageBar(result.data.UsedStorage, result.data.TotalStorage))
		h.LogResult.Printf(" \t\t%s\t(%s/%s)\n",
			result.data.StoreID,
			byteToVerboseString(result.data.UsedStorage),
			byteToVerboseString(result.data.TotalStorage),
		)
	}
	h.LogAbout.Println(getStorageBar(totalUsage, totalStorage))
	h.LogAbout.Println("Total Used Storage:", byteToVerboseString(totalUsage))
	h.LogAbout.Println("Total Storage:", byteToVerboseString(totalStorage))
}
