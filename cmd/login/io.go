package login

import (
	"fmt"
	"io"
	"syscall"

	"github.com/logrusorgru/aurora"
	"golang.org/x/term"
)

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

	return result
}

func successPrint(out io.Writer, name string) {
	fmt.Fprintf(out, "\n[%s] %s\n", name, aurora.Green("OK"))
}

func generatePrefix(username string) string {
	return "fssmtp"
}
