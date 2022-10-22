package credentials

const (
	CredentialsSecretName = "credentials"
	SMTPServerAddressKey  = "smtp-server-address"
	IMAPServerAddressKey  = "imap-server-address"
	UsernameKey           = "username"
	PasswordKey           = "password"
)

type Credentials struct {
	SMTPServerAddress string
	IMAPServerAddress string
	Username          string
	Password          string
}

type CredentialsStore interface {
	Put(string, map[string]string) error
	Get(string, string) (string, error)
}
