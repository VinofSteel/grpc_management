package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func passwordTagValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasSymbol              = regexp.MustCompile(`[^a-zA-Z0-9]`)
		hasUpperCasedCharacter = regexp.MustCompile(`[A-Z]`)
		hasLowerCasedCharacter = regexp.MustCompile(`[a-z]`)
		hasNumber              = regexp.MustCompile(`[0-9]`)
	)

	return hasSymbol.MatchString(password) && hasUpperCasedCharacter.MatchString(password) && hasLowerCasedCharacter.MatchString(password) && hasNumber.MatchString(password)
}
