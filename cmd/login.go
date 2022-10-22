package cmd

import (
	"github.com/deifyed/fsmail/cmd/login"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(0),
	RunE:  login.RunE(fs),
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
