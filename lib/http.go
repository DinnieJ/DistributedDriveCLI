package lib

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type Dictionary map[string]string

type RequestStruct struct {
	Url     string
	Method  string
	Headers Dictionary
	Query   Dictionary
	Body    []byte
}

const USER_AGENT = "DDCli/0.1.0"

var ErrFailedMetadata = errors.New("failed to get metadata")

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        20,
		TLSHandshakeTimeout: 0,
		DisableCompression:  true,
	},
}

func GetHttpRequest(h *RequestStruct) (*http.Request, error) {
	var r, err = http.NewRequest(
		h.Method,
		h.Url,
		If[io.Reader](len(h.Body) > 0, bytes.NewBuffer(h.Body), nil),
	)
	r.Header.Set("User-Agent", USER_AGENT)
	r.Header.Set("Accept", "*/*")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Cache-Control", "no-cache")
	if err != nil {
		return nil, err
	}
	if h.Headers != nil {
		for k, v := range h.Headers {
			r.Header.Set(k, v)
		}
	}
	var q = r.URL.Query()
	if h.Query != nil {
		for k, v := range h.Query {
			q.Add(k, v)
		}
	}
	r.URL.RawQuery = q.Encode()
	return r, nil
}

func HttpGetAccessToken(code string, clientId string, clientSecret string, codeVerifier string, redirect_uri string) (*http.Response, error) {
	var body = url.Values{}
	body.Set("code", code)
	body.Set("client_id", clientId)
	body.Set("client_secret", clientSecret)
	body.Set("code_verifier", codeVerifier)
	body.Set("redirect_uri", redirect_uri)
	body.Set("grant_type", "authorization_code")

	var request, err = GetHttpRequest(&RequestStruct{
		Url:     "https://oauth2.googleapis.com/token",
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:    []byte(body.Encode()),
	})
	if err != nil {
		return nil, err
	}

	return client.Do(request)
}

type UserInfoResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func HttpGetUserInfo(access_token string) *http.Response {
	var request, err = GetHttpRequest(&RequestStruct{
		Url:     "https://www.googleapis.com/oauth2/v1/userinfo",
		Method:  "GET",
		Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:    nil,
		Query:   map[string]string{"access_token": access_token},
	})
	if err != nil {
		return nil
	}

	var response = Must[*http.Response](client.Do(request))

	return response
}

func HttpRefreshAccessToken(clientId string, clientSecret string, refreshToken string) *http.Response {
	var body = url.Values{}
	body.Set("client_id", clientId)
	body.Set("client_secret", clientSecret)
	body.Set("refresh_token", refreshToken)
	body.Set("grant_type", "refresh_token")
	var request = Must[*http.Request](GetHttpRequest(&RequestStruct{
		Url:     "https://oauth2.googleapis.com/token",
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:    []byte(body.Encode()),
	}))

	return Must[*http.Response](client.Do(request))
}
