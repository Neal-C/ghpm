package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string
)

var rootCmd = &cobra.Command{
	Use:   "ghpm",
	Short: "ghpm is tool to manage privacy on github.",
	Args: cobra.ExactArgs(0),
	Long: `ghpm is tool to manage privacy on github. And quickly switch all repository to private.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if version != "" {
			return versionCmd.Execute()
		}

		fmt.Println("Made with 💞 love 💞 for developers by a developer ❤️")

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&version, "version", "v", "", "prints the version")
}
