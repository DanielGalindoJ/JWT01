package validator

import (
	"fmt"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("la contraseña debe tener al menos 6 caracteres")
	}
	return nil
}

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("el username debe tener al menos 3 caracteres")
	}
	if len(username) > 50 {
		return fmt.Errorf("el username no puede tener más de 50 caracteres")
	}
	return nil
}
