package errors

import (
	"fmt"
)

type AuthError struct {
	Err error
}

func (a AuthError) Error() string {
	return fmt.Sprintf("Authentication error: %v", a.Err)
}

func NewAuthError(inner error) AuthError {
	return AuthError{Err: inner}
}

type ValidationError struct {
	Err error
}

func (a ValidationError) Error() string {
	return fmt.Sprintf("Invalid resource state: %v", a.Err)
}

func NewValidationError(inner error) ValidationError {
	return ValidationError{Err: inner}
}

type RightsError struct {
	Err error
}

func (a RightsError) Error() string {
	return fmt.Sprintf("Unsufficient rights: %v", a.Err)
}

func NewRightsError(inner error) RightsError {
	return RightsError{Err: inner}
}
