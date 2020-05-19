package api

import (
	"fmt"
	"net/http"
)

// ErrResult reflects the body content when a request failed (g.e. Bad request or key conflict)
type ErrResult struct {
	Reason     string `json:"reason"`
	StatusCode int
}

func (err ErrResult) Error() string {
	return fmt.Sprintf("[%d]`%s`", err.StatusCode, err.Reason)
}

func (err ErrResult) IsTokenExpired() bool { return err.StatusCode == 412 }
func (err ErrResult) IsUnauthorized() bool { return err.StatusCode == http.StatusUnauthorized }
func (err ErrResult) IsForbidden() bool    { return err.StatusCode == http.StatusForbidden }
func (err ErrResult) IsNotFound() bool     { return err.StatusCode == http.StatusNotFound }
func (err ErrResult) IsRateLimited() bool  { return err.StatusCode == http.StatusTooManyRequests }
func (err ErrResult) IsSystemIssue() bool {
	return err.StatusCode >= 500 && err.StatusCode <= 599
}

// ErrAPIIncompatible reflect that the API version and the API server are misaligned
type ErrAPIIncompatible struct {
	Err          error
	DeveloperMsg string
}

func (err ErrAPIIncompatible) Error() string {
	return fmt.Sprintf("api incompatible: %v : %v", err.DeveloperMsg, err.Err)
}
func (err ErrAPIIncompatible) Unwrap() error { return err.Err }
