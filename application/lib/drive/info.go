package drive

import (
	"encoding/json"

	h "app.ddcli.datnn/application/lib"
)

func GetDriveInformation(driveConfig *h.Config) {
	if refreshToken, err := driveConfig.Get("refreshToken"); err == nil {
		h.LogResult.Println(refreshToken)
		var service = GetDriveService(refreshToken)
		var aboutResp, err = service.About.Get().Fields("*").Do()
		if err != nil {
			h.LogErr.Printf("Failed to get drive infomation: %s\n", err.Error())
		}
		h.LogResult.Printf("%v\n", string(h.Must[[]byte](json.Marshal(aboutResp))))
	}
}
