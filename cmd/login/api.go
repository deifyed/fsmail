package login

import (
	"fmt"
	"syscall"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func RunE(fs *afero.Afero) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, s []string) error {
		serverAddress := prompter("Server address: ", false)
		username := prompter("Username: ", false)
		password := prompter("Password: ", true)

		fmt.Println(serverAddress)
		fmt.Println(username)
		fmt.Println(password)

		return nil
	}
}

func prompter(msg string, hidden bool) string {
	fmt.Print(msg)

	var result string

	if hidden {
		rawResult, _ := term.ReadPassword(syscall.Stdin)

		result = string(rawResult)
	} else {
		_, err := fmt.Scanln(&result)
		if err != nil {
			panic(err.Error())
		}
	}

	fmt.Print("\n")

	return result
}
