// package for everything related to configuration
package config

var (
	// The "ghpm CLI" OAuth app
	// This value is safe to be embedded in version control
	OauthClientID = "Ov23li9W57aIXQHGJolI"

	// This value is safe to be embedded in version control
	// I know, you won't believe me, so : https://github.com/cli/cli/blob/trunk/internal/authflow/flow.go
	// feel free to shame github developers along with me ¯\_(ツ)_/¯
	OauthClientSecret = "9ff3b19b627a39639f1a0ffb2eb8e2416b5dcf02"

	CallbackURI = "http://127.0.0.1/callback"

	MinimumScopes = []string{"repo"}
)
