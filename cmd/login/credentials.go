package login

import (
	"fmt"

	"github.com/deifyed/fssmtp/pkg/credentials"
)

func validateCredentials(store credentials.CredentialsStore) error {
	_, err := store.Get(credentials.CredentialsSecretName, credentials.ServerAddressKey)
	if err != nil {
		return fmt.Errorf("retrieving server address: %w", err)
	}

	_, err = store.Get(credentials.CredentialsSecretName, credentials.UsernameKey)
	if err != nil {
		return fmt.Errorf("retrieving username: %w", err)
	}

	_, err = store.Get(credentials.CredentialsSecretName, credentials.PasswordKey)
	if err != nil {
		return fmt.Errorf("retrieving password: %w", err)
	}

	return nil
}

func promptForCredentials() credentials.Credentials {
	creds := credentials.Credentials{}

	creds.ServerAddress = prompter("Server address: ", false)
	creds.Username = prompter("Username: ", false)
	creds.Password = prompter("Password: ", true)

	return creds
}

func storeCredentials(store credentials.CredentialsStore, creds credentials.Credentials) error {
	err := store.Put(credentials.CredentialsSecretName, map[string]string{
		credentials.ServerAddressKey: creds.ServerAddress,
		credentials.UsernameKey:      creds.Username,
		credentials.PasswordKey:      creds.Password,
	})
	if err != nil {
		return fmt.Errorf("storing credentials: %w", err)
	}

	return nil
}
