package drive

import (
	"encoding/json"

	h "app.ddcli.datnn/lib"
	R "app.ddcli.datnn/root"
)

func GetAllDriveInfomation() {
	for k := range R.GoogleDriveCredential.Data {
		var config, _ = R.GoogleDriveCredential.GetConfig(k)
		if config == nil {
			continue
		}
		// if strings.HasPrefix(config.ConfigName, "google") {
		if refreshToken, err := config.Get("refreshToken"); err == nil {
			h.LogResult.Println(refreshToken)
			var service = GetDriveService(refreshToken)
			var aboutResp, err = service.About.Get().Fields("*").Do()
			if err != nil {
				h.LogErr.Printf("Failed to get drive infomation: %s\n", err.Error())
			}
			h.LogResult.Printf("%v\n", string(h.Must[[]byte](json.Marshal(aboutResp))))
		}
		// }
	}
}
