package keyring

import (
	"github.com/99designs/keyring"
)

func (c Client) open() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName:   c.Prefix,
		KeychainName:  c.Prefix,
		PassPrefix:    c.Prefix,
		WinCredPrefix: c.Prefix,
	})
}
