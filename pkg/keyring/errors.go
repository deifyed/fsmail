package keyring

import (
	"errors"
)

const (
	secretServiceItemNotFound = "The specified item could not be found in the keyring"
	secretServiceUserAborted  = "Cannot get secret of a locked object" //#nosec almost convinced this aint credentials
)

func handleError(err error, defaultError error) error {
	switch err.Error() {
	case secretServiceItemNotFound:
		return errors.New(err.Error())
	case secretServiceUserAborted:
		return errors.New(err.Error())
	default:
		return defaultError
	}
}
