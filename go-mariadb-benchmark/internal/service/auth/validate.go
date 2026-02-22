package auth

import (
	"errors"
	"fmt"
	"regexp"
)

func ValidateUsername(username string) error {
	regexUsername := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	if !regexUsername.MatchString(username) {
		return fmt.Errorf("validatePassword: %w", errors.New("username must only contain alphanumeric characters, _, -"))
	}
	if len(username) < 3 || len(username) > 32 {
		return fmt.Errorf("validateUsername: %w", errors.New("username must be between 3 and 32 characters long"))
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 128 {
		return fmt.Errorf("validatePassword: %w", errors.New("password must be between 8 and 128 characters long"))
	}
	regexPassword := regexp.MustCompile(`^[a-zA-Z0-9!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>/?\\|]+$`)
	if !regexPassword.MatchString(password) {
		return fmt.Errorf("validatePassword: %w", errors.New("password contains invalid character"))
	}
	return nil
}
