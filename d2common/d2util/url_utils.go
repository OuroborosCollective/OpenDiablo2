package d2util

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
)

// OpenURL opens the specified URL in the default browser, after validating it.
func OpenURL(rawURL string) error {
	if err := ValidateURL(rawURL); err != nil {
		return err
	}

	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", rawURL).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL).Start()
	case "darwin":
		return exec.Command("open", rawURL).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}

// ValidateURL checks if the URL is well-formed and uses an allowed scheme (http or https).
func ValidateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme: %s", u.Scheme)
	}

	return nil
}
