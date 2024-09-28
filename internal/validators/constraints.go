package validator

import (
	"regexp"
	"strings"
	"unicode"
)
type ValidationError struct {
	Key     string
	Message string
}

// Predefined validation errors
var (
	ErrNameTooShort        = ValidationError{"name", "name should be greater than %s characters"}
	ErrDescriptionTooShort = ValidationError{"description", "description is too short, must be at least 10 characters"}
	ErrFieldRequired       = ValidationError{"field", "this field cannot be empty"}
	ErrInvalidEmail        = ValidationError{"email", "invalid email address"}
	ErrPasswordTooWeak     = ValidationError{"password", "password is too weak, must include letters, numbers, and special characters"}
	ErrPhoneInvalid        = ValidationError{"phone", "invalid phone number"}
	ErrValueOutOfRange     = ValidationError{"value", "value is out of the acceptable range"}
	ErrInvalidDate         = ValidationError{"date", "invalid date format"}
)
func MinNameLength(name string, minLength int) bool {
	return len(strings.TrimSpace(name)) >= minLength
}

// MinDescriptionLength returns true if the description meets the minimum length.
func MinDescriptionLength(description string, minLength int) bool {
	return len(strings.TrimSpace(description)) >= minLength
}

// NotEmpty returns true if the field is not empty.
func NotEmpty(field string) bool {
	return len(strings.TrimSpace(field)) > 0
}

// IsValidEmail returns true if the email has a valid format.
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

// IsStrongPassword returns true if the password is strong (contains letters, numbers, and special characters).
func IsStrongPassword(password string) bool {
	var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool
	hasMinLen = len(password) >= 8
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

// IsValidPhone returns true if the phone number matches a simple pattern (e.g., US phone format).
func IsValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\+?(\d{1,3})?[-.\s]?\(?\d{1,4}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,9}$`)
	return re.MatchString(phone)
}

// IsInRange returns true if the value is within the given range.
func IsInRange(value, min, max int) bool {
	return value >= min && value <= max
}

// IsValidDate returns true if the date is in YYYY-MM-DD format.
func IsValidDate(date string) bool {
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return re.MatchString(date)
}
