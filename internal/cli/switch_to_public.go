package cli

import (
	"context"
	"log"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/Neal-C/ghpm/internal/ghpm"
	"github.com/spf13/cobra"
)

var switchToPublicCmd = &cobra.Command{
	Use:   "switch_public",
	Short: "Switch your private repository to public by name",
	Args:  cobra.ExactArgs(1),
	Long: heredoc.Docf(`
		Switch your public repository to private by name.

		Starts interactive setup and does a HTTP request to turn your repository public
	`, "`"),
	Example: heredoc.Doc(`
		# Starts interactive setup 
		and switches your repository to private by name
		
		$ ghpm switch_private
		`),
	RunE: func(cmd *cobra.Command, args []string) error {

		token, err := ghpm.LoginToGithubWithDetecFlow()

		if err != nil {
			return err
		}

		ghPrivacyManager := ghpm.NewGithubPrivacyManager(token, http.DefaultClient)

		name := args[0]

		err = ghPrivacyManager.SwitchRepoToPublicByName(context.Background(), name)

		if err != nil {
			return err
		}

		log.Printf("success. %s was made public", name)

		return nil

	},
}

func init() {
	rootCmd.AddCommand(switchToPublicCmd)
}
