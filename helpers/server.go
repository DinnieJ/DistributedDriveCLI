package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

const DEFAULT_PORT = 28899

const HTML_TEMPLATE = `
<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>DDCLI</title>
		</head>
		<body>
			<h1>
				%s
			</h1>
			<p>%s</p>
			<pre>%v</pre>
		</body>
	</html>
`

func ResponseError(w io.Writer, err string, err_detail interface{}) {
	fmt.Fprintf(w, HTML_TEMPLATE, "Failed to add Google Drive to application", err, err_detail)
}

func ResponseSuccess(w io.Writer) {
	fmt.Fprintf(w, HTML_TEMPLATE, "Added New Drive to application successfully", "You can close this browser", nil)
}

type Credential struct {
	mu           sync.Mutex
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (c *Credential) SetToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.AccessToken = token
}

func (c *Credential) GetToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.AccessToken
}

func StartCredentialCallbackServer(port int, appConfig *Configuration) (*Credential, error) {
	var listenerPort = DEFAULT_PORT
	if port != 0 {
		listenerPort = port
	}
	var oauthCodeChallenge = GenerateCodeChallengeS256()
	var clientId = Must[string](appConfig.GetOrError("googleAppClient", "clientID"))

	var oauthUrl = "https://accounts.google.com/o/oauth2/v2/auth?" +
		(fmt.Sprintf(`scope=%s&response_type=code&client_id=%s&prompt=select_account&redirect_uri=%s&access_type=offline&code_challenge=%s&code_challenge_method=S256`,
			strings.Join([]string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/drive",
			}, " ",
			),
			clientId,
			fmt.Sprintf("http://127.0.0.1:%d/callback", listenerPort),
			oauthCodeChallenge.CodeChallenge,
		))
	openBrowserWithUrl(oauthUrl)

	var mux = http.NewServeMux()
	mux.HandleFunc("/callback", CallbackWrapper)

	var ctx, cancel = context.WithCancel(context.Background())
	var result = &Credential{}

	var addr = Spr("127.0.0.1:%d", listenerPort)
	LogInfo.Printf("Start callback server at %s\n", addr)
	var server = &http.Server{
		Addr:    addr,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, "appConfig", appConfig)
			ctx = context.WithValue(ctx, "cancel", cancel)
			ctx = context.WithValue(ctx, "oauthCodeChallenge", oauthCodeChallenge)
			ctx = context.WithValue(ctx, "result", result)
			return ctx
		},
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-ctx.Done() // Waiting for finish handle callback request
	if err := server.Shutdown(context.Background()); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return nil, err
	}
	return result, nil
}

func CallbackWrapper(resp http.ResponseWriter, request *http.Request) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	var ctx = request.Context()
	var cancel = ctx.Value("cancel").(context.CancelFunc)
	defer cancel()
	var result = ctx.Value("result").(*Credential)
	var appConfig = ctx.Value("appConfig").(*Configuration)
	var oauthCodeChallenge = ctx.Value("oauthCodeChallenge").(*OAuthCodeChallenge)
	if err := request.URL.Query().Get("error"); err != "" {
		ResponseError(resp, err, nil)
	}
	if code := request.URL.Query().Get("code"); code != "" {
		var clientId = Must[string](appConfig.GetOrError("googleAppClient", "clientID"))
		var clientSecret = Must[string](appConfig.GetOrError("googleAppClient", "clientSecret"))
		var access_token_response = Must[*http.Response](HttpGetAccessToken(
			code,
			clientId,
			clientSecret,
			oauthCodeChallenge.CodeVerifier,
			fmt.Sprintf("http://%s%s", request.Host, request.URL.Path),
		))
		defer access_token_response.Body.Close()

		var body = Must[[]byte](io.ReadAll(access_token_response.Body))
		result.mu.Lock()
		defer result.mu.Unlock()
		if err := json.Unmarshal(body, &result); err != nil {
			LogErr.Printf("Failed to get access token: %s\n", err.Error())
			ResponseError(resp, "Failed to get access token", string(body))
		} else {
			ResponseSuccess(resp)
		}
	}
}

func openBrowserWithUrl(url string) {
	var err error
	LogResult.Printf("Access URL %s to login to Google\n", url)
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		LogErr.Println(err.Error())
		// LogResult.Printf("Failed to open browser: Access URL %s to login to Google", url)
	}
}
