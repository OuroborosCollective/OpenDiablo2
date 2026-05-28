package d2server

import "testing"

func TestGetBindAddr(t *testing.T) {
	tests := []struct {
		networkServer bool
		want          string
	}{
		{false, "127.0.0.1"},
		{true, "0.0.0.0"},
	}

	for _, tt := range tests {
		g := &GameServer{
			networkServer: tt.networkServer,
		}
		if got := g.getBindAddr(); got != tt.want {
			t.Errorf("getBindAddr() = %v, want %v", got, tt.want)
		}
	}
}
