package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)
type ValidationError struct {
	Key     string
	Message string
}
var(
	MIN_TASK_NAME_LENGTH=3
	MAX_TASK_NAME_LENGTH=40
	MIN_TASK_DESCRIPTION_LENGTH=10
	MAX_TASK_DESCRIPTION_LENGTH=200
)
var (
	ErrNameTooShort        = ValidationError{"name", fmt.Sprintf("name should be between %d - %d characters", MIN_TASK_NAME_LENGTH, MAX_TASK_NAME_LENGTH)}
    ErrDescriptionTooShort = ValidationError{"description", fmt.Sprintf("description must be between %d - %d characters", MIN_TASK_DESCRIPTION_LENGTH, MAX_TASK_DESCRIPTION_LENGTH)}
	ErrFieldRequired       = ValidationError{"field", "this field cannot be empty"}
	ErrInvalidEmail        = ValidationError{"email", "invalid email address"}
	ErrPasswordTooWeak     = ValidationError{"password", "password is too weak, must include letters, numbers, and special characters"}
	ErrPhoneInvalid        = ValidationError{"phone", "invalid phone number"}
	ErrValueOutOfRange     = ValidationError{"value", "value is out of the acceptable range"}
	ErrInvalidDate         = ValidationError{"date", "invalid date format"}
)
func MinNameLength(name string) bool {
    return len(strings.TrimSpace(name)) >= MIN_TASK_NAME_LENGTH && len(name) < MAX_TASK_NAME_LENGTH
}

func MinDescriptionLength(description string) bool {
    return len(strings.TrimSpace(description)) >= MIN_TASK_DESCRIPTION_LENGTH && len(description) < MAX_TASK_DESCRIPTION_LENGTH
}

func NotEmpty(field string) bool {
	return len(strings.TrimSpace(field)) > 0
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

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

func IsValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\+?(\d{1,3})?[-.\s]?\(?\d{1,4}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,9}$`)
	return re.MatchString(phone)
}

func IsInRange(value, min, max int) bool {
	return value >= min && value <= max
}

func IsValidDate(date string) bool {
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return re.MatchString(date)
}
