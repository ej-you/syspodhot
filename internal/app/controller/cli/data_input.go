package cli

import (
	"errors"
	"fmt"

	"syspodhot/internal/pkg/password"
)

const _usernameMaxLen = 50 // max length of username

// userBody represents a user data.
type userBody struct {
	// user username
	Username string
	// user password
	Password []byte
}

// ParseUsername parses user username from given string and validates it.
func (u *userBody) ParseUsername(raw string) error {
	if len(raw) > _usernameMaxLen {
		return fmt.Errorf(
			"username is too long: accepted %d (or less) simbols length",
			_usernameMaxLen,
		)
	}
	u.Username = raw
	return nil
}

// ParsePassword parses user password from given string and validates it.
func (u *userBody) ParsePassword(raw string) error {
	if !password.Strong(raw) {
		return errors.New("weak password")
	}
	u.Password = []byte(raw)
	return nil
}
