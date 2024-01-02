package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

const DEFAULT_PORT = 28899

type Credential struct {
	mu    sync.Mutex
	Token string
}

func (c *Credential) SetToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Token = token
}

func (c *Credential) GetToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Token
}

func StartCredentialCallbackServer(port int, appConfig *Configuration) (*Credential, error) {
	var listenerPort = DEFAULT_PORT
	if port != 0 {
		listenerPort = port
	}
	var oauthCodeChallenge = GenerateCodeChallengeS256()
	var clientId, err = appConfig.GetOrError("googleAppClient", "clientID")
	if err != nil {
		return nil, err
	}

	var oauthUrl = "https://accounts.google.com/o/oauth2/v2/auth?" +
		(fmt.Sprintf(`scope=%s&response_type=code&client_id=%s&redirect_uri=%s&prompt=consent&access_type=offline&code_challenge=%s&code_challenge_method=S256`,
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
	var ctx = request.Context()
	var cancel = ctx.Value("cancel").(context.CancelFunc)
	var result = ctx.Value("result").(*Credential)
	var code = request.URL.Query().Get("code")
	if code != "" {
		resp.Header().Set("Content-Type", "text/html; charset=utf-8")
		const successHtml = `
		<!DOCTYPE html>
			<html lang="en">
				<head>
					<meta charset="UTF-8">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
					<title>DDCLI</title>
				</head>
				<body>
					<h1>
						Add new Drive successfully
					</h1>
					<p>Now you can close this windows safely</p>
				</body>
			</html>
		`
		fmt.Fprint(resp, successHtml)
		result.SetToken(code)
		cancel()
	}
}

func openBrowserWithUrl(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Unsupported platform")
	}
	if err != nil {
		LogErr.Println(err.Error())
		LogResult.Printf("Failed to open browser: Access URL %s to login to Google", url)
	}
}
