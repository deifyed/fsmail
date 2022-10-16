package cmd

import (
	"github.com/deifyed/fssmtp/cmd/login"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	RunE:  login.RunE(fs),
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
