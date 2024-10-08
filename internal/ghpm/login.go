package ghpm

import (
	"net/http"

	"github.com/Neal-C/ghpm/internal/config"
	"github.com/cli/oauth"
)

func LoginToGithubWithDetecFlow() (string, error) {
	flow := &oauth.Flow{
		Host:         oauth.GitHubHost("https://github.com"),
		ClientID:     config.OauthClientID,
		ClientSecret: config.OauthClientSecret,
		CallbackURI:  config.CallbackURI,
		Scopes:       config.MinimumScopes,
		HTTPClient:   http.DefaultClient,
	}

	accessToken, err := flow.DetectFlow()

	if err != nil {
		return "", err
	}

	return accessToken.Token, err
}
