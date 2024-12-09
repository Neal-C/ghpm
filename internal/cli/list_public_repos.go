package cli

import (
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/Neal-C/ghpm/internal/ghpm"
	"github.com/spf13/cobra"
)

var listAllPublicRepositoriesCmd = &cobra.Command{
	Use:   "list_public",
	Short: "List all your public repositories.",
	Args:  cobra.NoArgs,
	Long: heredoc.Docf(`
		List all your public repositories.
	`, "`"),
	Example: heredoc.Doc(`
		# Starts interactive setup 
		and lists your public repositories

		Only the first 100 in alphabetical order.
		
		$ ghpm list_public
		`),
	RunE: func(cmd *cobra.Command, args []string) error {

		token, err := ghpm.LoginToGithubWithDetecFlow()

		if err != nil {
			return err
		}

		ghPrivacyManager := ghpm.NewGithubPrivacyManager(token, http.DefaultClient)

		err = ghPrivacyManager.ListAllPublicRepositories(cmd.Context())

		if err != nil {
			return err
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(listAllPublicRepositoriesCmd)
}
