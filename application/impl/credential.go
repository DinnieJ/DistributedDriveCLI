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

	R.GoogleDriveCredential.SetConfig(h.Spr("GoogleDrive|%s", userInfoResponse.Id), "refreshToken", credential.RefreshToken)
	R.GoogleDriveCredential.SetConfig(h.Spr("GoogleDrive|%s", userInfoResponse.Id), "email", userInfoResponse.Email)
	R.GoogleDriveCredential.SaveConfig()
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

func getStorageBar(value float64, rng float64) string {
	var storagePercent = value / rng
	var barLength = 30
	var usedBarLength = math.Ceil(storagePercent * float64(barLength))
	var bar = ""
	bar += ""
	bar += strings.Repeat("█", int(usedBarLength))
	bar += strings.Repeat("░", barLength-int(usedBarLength))
	bar += ""
	bar += fmt.Sprintf(" %.2f %%", storagePercent*100)
	return bar
}

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

func GetAllStorageInformation() {
	var wg = &sync.WaitGroup{}
	var mutex = &sync.Mutex{}
	var lStore = len(R.GoogleDriveCredential.Data)
	wg.Add(lStore)
	var listResult = make([]resultStorageInformation, 0, lStore)
	for k := range R.GoogleDriveCredential.Data {
		var config, err = R.GoogleDriveCredential.GetConfig(k)
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
		h.LogResult.Printf("%s \t\t%s\t(%s/%s)\n",
			getStorageBar(result.data.UsedStorage, result.data.TotalStorage),
			result.data.StoreID,
			byteToVerboseString(result.data.UsedStorage),
			byteToVerboseString(result.data.TotalStorage),
		)
	}
	h.LogAbout.Println(getStorageBar(totalUsage, totalStorage))
	h.LogAbout.Println("Total Used Storage:", byteToVerboseString(totalUsage))
	h.LogAbout.Println("Total Storage:", byteToVerboseString(totalStorage))
}
