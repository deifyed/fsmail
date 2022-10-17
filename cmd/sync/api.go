package sync

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/deifyed/fssmtp/pkg/credentials"
	"github.com/deifyed/fssmtp/pkg/keyring"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/gomail.v2"
)

func RunE(fs *afero.Afero) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		store := keyring.Client{Prefix: generatePrefix("")}

		serverAddress, err := store.Get(credentials.CredentialsSecretName, credentials.ServerAddressKey)
		if err != nil {
			return fmt.Errorf("retrieving server address: %w", err)
		}

		username, err := store.Get(credentials.CredentialsSecretName, credentials.UsernameKey)
		if err != nil {
			return fmt.Errorf("retrieving username: %w", err)
		}

		password, err := store.Get(credentials.CredentialsSecretName, credentials.PasswordKey)
		if err != nil {
			return fmt.Errorf("retrieving password: %w", err)
		}

		parts := strings.Split(serverAddress, ":")
		host := parts[0]

		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("converting port from string to int: %w", err)
		}

		dialer := gomail.NewDialer(host, port, username, password)

		sender, err := dialer.Dial()
		if err != nil {
			return fmt.Errorf("dialing: %w", err)
		}

		defer func() {
			_ = sender.Close()
		}()

		return nil
	}
}

func generatePrefix(username string) string {
	return "fssmtp"
}
