package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to github.",
	Args:  cobra.NoArgs,
	Long: heredoc.Docf(`
		Authenticate with a GitHub host.

		The default hostname is %[1]sgithub.com%[1]s. This can be overridden using the %[1]s--hostname%[1]s
		flag.

		The default authentication mode is a web-based browser flow. After completion, an
		authentication token will be stored securely in the system credential store.
		If a credential store is not found or there is an issue using it ghpm will fallback
		to writing the token to a plain text file. See %[1]sgh auth status%[1]s for its
		stored location.

		Alternatively, use %[1]s--with-token%[1]s to pass in a token on standard input.
		The minimum required scopes for the token are: %[1]srepo%[1]s, %[1]sread:org%[1]s, and %[1]sgist%[1]s.

		Alternatively, ghpm will use the authentication token found in environment variables.
		This method is most suitable for "headless" use of gh such as in automation. See
		%[1]sgh help environment%[1]s for more info.

		The git protocol to use for git operations on this host can be set with %[1]s--git-protocol%[1]s,
		or during the interactive prompting. Although login is for a single account on a host, setting
		the git protocol will take effect for all users on the host.

		Specifying %[1]sssh%[1]s for the git protocol will detect existing SSH keys to upload,
		prompting to create and upload a new key if one is not found. This can be skipped with
		%[1]s--skip-ssh-key%[1]s flag.
	`, "`"),
	Example: heredoc.Doc(`
		# Start interactive setup
		$ ghpm login

		# Authenticate against github.com by reading the token from a file
		$ ghpm login --with-token < mytoken.txt
		`),
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("this command does nothing for now. Placeholder in case I want to implement persist auth across multiple commands")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
