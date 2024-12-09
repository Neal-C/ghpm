package cli

import (
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/Neal-C/ghpm/internal/ghpm"
	"github.com/spf13/cobra"
)

var switchAllToPrivateCmd = &cobra.Command{
	Use:   "thanos_snap",
	Short: "Switch all your public repositories to private.",
	Args:  cobra.NoArgs,
	Long: heredoc.Docf(`
		Switch all your public repositories to private.

		By default, starred repositories with 1 stars are not turned private.

		Starts interactive setup and does a HTTP request against all your public repositories to turn them private
	`, "`"),
	Example: heredoc.Doc(`
		# Starts interactive setup 
		and request all your public repositories to turn private
		
		$ ghpm thanos_snap
		`),
	RunE: func(cmd *cobra.Command, args []string) error {

		token, err := ghpm.LoginToGithubWithDetecFlow()

		if err != nil {
			return err
		}

		ghPrivacyManager := ghpm.NewGithubPrivacyManager(token, http.DefaultClient)

		err = ghPrivacyManager.SwitchAllRepositoriesToPrivate(cmd.Context())

		if err != nil {
			return err
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(switchAllToPrivateCmd)
}
