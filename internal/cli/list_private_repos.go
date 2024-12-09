package cli

import (
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/Neal-C/ghpm/internal/ghpm"
	"github.com/spf13/cobra"
)

var listAllPrivateRepositoriesCmd = &cobra.Command{
	Use:   "list_private",
	Short: "List all your private repositories.",
	Args:  cobra.NoArgs,
	Long: heredoc.Docf(`
		List all your private repositories.
	`, "`"),
	Example: heredoc.Doc(`
		# Starts interactive setup 
		and lists your private repositories

		Only the first 100 in alphabetical order.
		
		$ ghpm list_private
		`),
	RunE: func(cmd *cobra.Command, args []string) error {

		token, err := ghpm.LoginToGithubWithDetecFlow()

		if err != nil {
			return err
		}

		ghPrivacyManager := ghpm.NewGithubPrivacyManager(token, http.DefaultClient)

		err = ghPrivacyManager.ListAllPrivateRepositories(cmd.Context())

		if err != nil {
			return err
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(listAllPrivateRepositoriesCmd)
}
