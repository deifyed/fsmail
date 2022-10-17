package login

import "fmt"

const (
	credentialsSecretName = "credentials"
	serverAddressKey      = "server-address"
	usernameKey           = "username"
	passwordKey           = "password"
)

type credentials struct {
	serverAddress string
	username      string
	password      string
}

type credentialsStore interface {
	Put(string, map[string]string) error
	Get(string, string) (string, error)
}

func validateCredentials(store credentialsStore) error {
	_, err := store.Get(credentialsSecretName, serverAddressKey)
	if err != nil {
		return fmt.Errorf("retrieving server address: %w", err)
	}

	_, err = store.Get(credentialsSecretName, usernameKey)
	if err != nil {
		return fmt.Errorf("retrieving username: %w", err)
	}

	_, err = store.Get(credentialsSecretName, passwordKey)
	if err != nil {
		return fmt.Errorf("retrieving password: %w", err)
	}

	return nil
}

func promptForCredentials() credentials {
	creds := credentials{}

	creds.serverAddress = prompter("Server address: ", false)
	creds.username = prompter("Username: ", false)
	creds.password = prompter("Password: ", true)

	return creds
}

func storeCredentials(store credentialsStore, creds credentials) error {
	err := store.Put(credentialsSecretName, map[string]string{
		serverAddressKey: creds.serverAddress,
		usernameKey:      creds.username,
		passwordKey:      creds.password,
	})
	if err != nil {
		return fmt.Errorf("storing credentials: %w", err)
	}

	return nil
}
