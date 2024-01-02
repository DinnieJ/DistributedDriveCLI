package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

type OAuthCodeChallenge struct {
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

func GenerateCodeChallengeS256() *OAuthCodeChallenge {
	var codeVerifer = base64.RawURLEncoding.EncodeToString(randomBytes(32))
	var s256CodeVerifier = sha256.Sum256([]byte(codeVerifer))
	var CodeChallenge = base64.RawURLEncoding.EncodeToString([]byte(s256CodeVerifier[:]))
	CodeChallenge = strings.Replace(CodeChallenge, "+", "-", -1)
	CodeChallenge = strings.Replace(CodeChallenge, "/", "_", -1)
	CodeChallenge = strings.Replace(CodeChallenge, "=", "", -1)
	return &OAuthCodeChallenge{
		CodeVerifier:        string(codeVerifer),
		CodeChallenge:       CodeChallenge,
		CodeChallengeMethod: "S256",
	}
}
