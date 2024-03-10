package middleware

import (
	"errors"
)

func CheckAdmin(isAdmin bool) error {
	if !isAdmin {
		return errors.New("acced denied")
	}

	return nil
}
