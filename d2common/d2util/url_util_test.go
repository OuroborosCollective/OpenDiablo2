package d2util

import "testing"

func TestIsValidBrowserURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://www.github.com/OpenDiablo2/OpenDiablo2", true},
		{"http://example.com", true},
		{"ftp://example.com", false},
		{"file:///etc/passwd", false},
		{"javascript:alert(1)", false},
		{"data:text/plain,hello", false},
		{"mailto:test@example.com", false},
		{"", false},
		{"not-a-url", false},
		{"http://", true}, // url.Parse might return true for this, but scheme is http
		{"https://", true},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			if got := IsValidBrowserURL(tt.url); got != tt.expected {
				t.Errorf("IsValidBrowserURL(%q) = %v, want %v", tt.url, got, tt.expected)
			}
		})
	}
}
