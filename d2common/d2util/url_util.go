package d2util

import (
	"net/url"
)

// IsValidBrowserURL checks if the given string is a valid http or https URL.
func IsValidBrowserURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return u.Scheme == "http" || u.Scheme == "https"
}
