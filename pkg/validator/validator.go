package validator

import (
	"fmt"
	"regexp"
)

type ValidationError struct {
	Field   string
	Message string
}

type Validator struct {
	errors []ValidationError
}

func New() *Validator {
	return &Validator{
		errors: []ValidationError{},
	}
}

func (v *Validator) ValidateEmail(email string) *Validator {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(email) {
		v.errors = append(v.errors, ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		})
	}
	return v
}

func (v *Validator) ValidatePassword(password string) *Validator {
	if len(password) < 8 {
		v.errors = append(v.errors, ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
		})
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		v.errors = append(v.errors, ValidationError{
			Field:   "password",
			Message: "Password must contain at least one lowercase letter",
		})
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		v.errors = append(v.errors, ValidationError{
			Field:   "password",
			Message: "Password must contain at least one uppercase letter",
		})
	}

	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		v.errors = append(v.errors, ValidationError{
			Field:   "password",
			Message: "Password must contain at least one number",
		})
	}

	return v
}

func (v *Validator) ValidateRequired(field, value string) *Validator {
	if value == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s is required", field),
		})
	}
	return v
}

func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

func (v *Validator) GetErrors() []ValidationError {
	return v.errors
}
