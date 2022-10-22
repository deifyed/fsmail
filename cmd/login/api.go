package login

import (
	"fmt"

	"github.com/deifyed/fsmail/pkg/keyring"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunE(fs *afero.Afero) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		creds := promptForCredentials()

		keyringClient := keyring.Client{Prefix: generatePrefix(creds.Username)}

		err := storeCredentials(keyringClient, creds)
		if err != nil {
			return fmt.Errorf("storing credentials: %w", err)
		}

		err = validateCredentials(keyringClient)
		if err != nil {
			return fmt.Errorf("validating credentials: %w", err)
		}

		successPrint(cmd.OutOrStdout(), "Credentials")

		return nil
	}
}
