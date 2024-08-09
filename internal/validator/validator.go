package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Define a Validator struct which contains a
// map of validation error messages for form fields
type Validator struct {
	FieldErrors map[string]string
}

// Return true is FieldErrors has no entries
// i.e., there are no errors.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldErrors() adds error message to the FieldErrors map
// (as long as no entry already exists for the given key)
func (v *Validator) AddFieldErrors(key, message string) {

	// Initialize map first if it doesn't exist
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField() adds an error message to the FieldErrors Map
// only if a validation check is NOT okay
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldErrors(key, message)
	}
}

// NotBlank() returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if the value is less than or equal to n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedValued[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// MinChars() returns true if value contains at least n chars
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
