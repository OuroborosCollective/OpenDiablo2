package d2util

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		url     string
		wantErr bool
	}{
		{"https://www.github.com/OpenDiablo2/OpenDiablo2", false},
		{"http://example.com", false},
		{"ftp://example.com", true},
		{"javascript:alert(1)", true},
		{"file:///etc/passwd", true},
		{"invalid-url", true},
		{"https://github.com/OpenDiablo2/OpenDiablo2;calc.exe", false}, // scheme is still https, url.Parse might allow it depending on how it's used, but ValidateURL only checks scheme and Parse error.
	}

	for _, tt := range tests {
		err := ValidateURL(tt.url)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
		}
	}
}
