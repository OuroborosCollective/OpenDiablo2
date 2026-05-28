package d2server

import (
	"net/http"
	"testing"
)

func TestCheckOrigin(t *testing.T) {
	g := &GameServer{}

	tests := []struct {
		name           string
		origin         string
		host           string
		allowedOrigins []string
		want           bool
	}{
		{
			name:   "Empty Origin",
			origin: "",
			host:   "localhost:6670",
			want:   true,
		},
		{
			name:   "Matching Host",
			origin: "http://localhost:3000",
			host:   "localhost:6670",
			want:   true,
		},
		{
			name:   "Matching Host IP",
			origin: "http://127.0.0.1:3000",
			host:   "127.0.0.1:6670",
			want:   true,
		},
		{
			name:   "Disallowed Host",
			origin: "http://malicious.com",
			host:   "localhost:6670",
			want:   false,
		},
		{
			name:           "Allowed via list (hostname)",
			origin:         "http://trusted.com:8080",
			host:           "localhost:6670",
			allowedOrigins: []string{"trusted.com"},
			want:           true,
		},
		{
			name:           "Allowed via list (full origin)",
			origin:         "http://another.com",
			host:           "localhost:6670",
			allowedOrigins: []string{"http://another.com"},
			want:           true,
		},
		{
			name:   "Case Insensitive Host Match",
			origin: "http://LOCALHOST:3000",
			host:   "localhost:6670",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.SetAllowedOrigins(tt.allowedOrigins)
			r, _ := http.NewRequest("GET", "ws://"+tt.host+"/ws", nil)
			if tt.origin != "" {
				r.Header.Set("Origin", tt.origin)
			}
			if got := g.checkOrigin(r); got != tt.want {
				t.Errorf("checkOrigin() for %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
