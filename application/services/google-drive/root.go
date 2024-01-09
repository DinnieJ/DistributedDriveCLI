package drive

import (
	"context"
	"encoding/json"
	"io"

	R "app.ddcli.datnn/application"
	h "app.ddcli.datnn/application/lib"
	s "app.ddcli.datnn/application/services"
	"golang.org/x/oauth2"
	d "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type customTokenResource struct {
	RefreshToken string
}

func (c *customTokenResource) Token() (*oauth2.Token, error) {
	var token = &oauth2.Token{}
	var response = h.HttpRefreshAccessToken(
		h.Must[string](R.AppConfiguration.GetOrError("googleAppClient", "clientID")),
		h.Must[string](R.AppConfiguration.GetOrError("googleAppClient", "clientSecret")),
		c.RefreshToken,
	)
	defer response.Body.Close()

	var bbody = h.Must(io.ReadAll(response.Body))
	if err := json.Unmarshal(bbody, &token); err != nil {
		h.LogErr.Printf("Failed to get access token: %s\n", err.Error())
		return nil, err
	}

	return token, nil
}

type GoogleDriveService struct {
	driveService *d.Service
	config       *h.Config
}

func (g *GoogleDriveService) New(config *h.Config) error {
	var refreshToken = h.Must(config.Get("refreshToken"))
	var customTokenResource = &customTokenResource{
		RefreshToken: refreshToken,
	}
	var service, err = d.NewService(context.TODO(), option.WithTokenSource(customTokenResource))
	if err != nil {
		return err
	}
	g.driveService = service
	g.config = config
	return nil
}

func (g *GoogleDriveService) GetInformation() (*s.StorageInformation, error) {
	var aboutResp, err = g.driveService.About.Get().Fields("*").Do()
	return &s.StorageInformation{
		StoreID:      h.Must(g.config.Get("email")),
		StoreName:    "google",
		RootFolder:   "",
		TotalStorage: float64(aboutResp.StorageQuota.Limit),
		UsedStorage:  float64(aboutResp.StorageQuota.Usage),
	}, err
}
