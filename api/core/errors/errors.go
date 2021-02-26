package errors

import "fmt"

type AuthError struct {
	Err error
}

func (a AuthError) Error() string {
	return fmt.Sprintf("Authentication error: %v", a.Err)
}

