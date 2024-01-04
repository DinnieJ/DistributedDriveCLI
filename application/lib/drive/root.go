package drive

import (
	"context"
	"encoding/json"
	"io"

	R "app.ddcli.datnn/application"
	h "app.ddcli.datnn/application/lib"
	"golang.org/x/oauth2"
	d "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type CustomTokenSource struct {
	RefreshToken string
}

func (c *CustomTokenSource) Token() (*oauth2.Token, error) {
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

func GetDriveService(refreshToken string) *d.Service {
	var GetOAuth2Token = &CustomTokenSource{RefreshToken: refreshToken}
	var service, err = d.NewService(context.TODO(), option.WithTokenSource(GetOAuth2Token))
	if err != nil {
		h.LogErr.Printf("Failed to get drive service: %s\n", err.Error())
	}
	return service
}
