package cli

import (
	"context"
	"log"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/Neal-C/ghpm/internal/ghpm"
	"github.com/spf13/cobra"
)

var switchToPrivateCmd = &cobra.Command{
	Use:   "switch_private",
	Short: "Switch your public repository to private by name",
	Args:  cobra.ExactArgs(1),
	Long: heredoc.Docf(`
		Switch your public repository to private by name.

		By default, starred repositories with 4 will be made private.

		Starts interactive setup and does a HTTP request to turn your repository private.
	`, "`"),
	Example: heredoc.Doc(`
		# Starts interactive setup 
		and switches your repository to private by name
		
		$ ghpm switch_public <name here>
		`),
	RunE: func(cmd *cobra.Command, args []string) error {

		token, err := ghpm.LoginToGithubWithDetecFlow()

		if err != nil {
			return err
		}

		ghPrivacyManager := ghpm.NewGithubPrivacyManager(token, http.DefaultClient)

		name := args[0]

		err = ghPrivacyManager.SwitchRepoToPrivateByName(context.Background(), name)

		if err != nil {
			return err
		}

		log.Printf("success. %s was made private",name)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchToPrivateCmd)
}
