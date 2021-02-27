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


type ExistingResourceError struct {
	Err error
}

func (a ExistingResourceError) Error() string {
	return fmt.Sprintf("Resource already exists: %v", a.Err)
}

func NewExistingResourceError(inner error) ExistingResourceError {
	return ExistingResourceError{Err: inner}
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
