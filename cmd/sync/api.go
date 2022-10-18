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
		creds, err := acquireCredentials()
		if err != nil {
			return fmt.Errorf("acquiring credentials: %w", err)
		}

		parts := strings.Split(creds.ServerAddress, ":")
		host := parts[0]

		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("converting port from string to int: %w", err)
		}

		dialer := gomail.NewDialer(host, port, creds.Username, creds.Password)

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

func acquireCredentials() (credentials.Credentials, error) {
	store := keyring.Client{Prefix: generatePrefix("")}

	creds := credentials.Credentials{}
	var err error

	creds.ServerAddress, err = store.Get(credentials.CredentialsSecretName, credentials.ServerAddressKey)
	if err != nil {
		return credentials.Credentials{}, fmt.Errorf("retrieving server address: %w", err)
	}

	creds.Username, err = store.Get(credentials.CredentialsSecretName, credentials.UsernameKey)
	if err != nil {
		return credentials.Credentials{}, fmt.Errorf("retrieving username: %w", err)
	}

	creds.Password, err = store.Get(credentials.CredentialsSecretName, credentials.PasswordKey)
	if err != nil {
		return credentials.Credentials{}, fmt.Errorf("retrieving password: %w", err)
	}

	return creds, nil
}

func generatePrefix(username string) string {
	return "fssmtp"
}
