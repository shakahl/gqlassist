package utils

import (
	"net/url"
)

// IsValidURL checks if a string is valid url
func IsValidURL(u string) bool {
	if u == "" {
		return false
	}
	if _, err := url.ParseRequestURI(u); err != nil {
		return false
	}
	return true
}
