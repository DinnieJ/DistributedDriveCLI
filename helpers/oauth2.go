package helpers

type OAuthCodeChallenge struct {
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

func (o *OAuthCodeChallenge) GenerateCodeChallengeS256 {
	return
}
