package credentials

const (
	CredentialsSecretName = "credentials"
	ServerAddressKey      = "server-address"
	UsernameKey           = "username"
	PasswordKey           = "password"
)

type Credentials struct {
	ServerAddress string
	Username      string
	Password      string
}

type CredentialsStore interface {
	Put(string, map[string]string) error
	Get(string, string) (string, error)
}
