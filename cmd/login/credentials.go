package login

import (
	"fmt"

	"github.com/deifyed/fsmail/pkg/credentials"
)

func validateCredentials(store credentials.CredentialsStore) error {
	_, err := store.Get(credentials.CredentialsSecretName, credentials.SMTPServerAddressKey)
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

	creds.SMTPServerAddress = prompter("SMTP server address: ", false)
	creds.IMAPServerAddress = prompter("IMAP server address: ", false)
	creds.Username = prompter("Username: ", false)
	creds.Password = prompter("Password: ", true)

	return creds
}

func storeCredentials(store credentials.CredentialsStore, creds credentials.Credentials) error {
	err := store.Put(credentials.CredentialsSecretName, map[string]string{
		credentials.SMTPServerAddressKey: creds.SMTPServerAddress,
		credentials.IMAPServerAddressKey: creds.IMAPServerAddress,
		credentials.UsernameKey:          creds.Username,
		credentials.PasswordKey:          creds.Password,
	})
	if err != nil {
		return fmt.Errorf("storing credentials: %w", err)
	}

	return nil
}
