package d2server

import (
	"fmt"
	"net/http"
)

// StartWebSocket starts the WebSocket server on the given port
func (g *GameServer) StartWebSocket(wsPort int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", g.HandleWebSocket)

	bindAddr := g.getBindAddr()

	g.Infof("Starting WebSocket Game Server @ %s:%d/ws\n", bindAddr, wsPort)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", bindAddr, wsPort), mux)
		if err != nil {
			g.Errorf("WebSocket server failed: %v", err)
		}
	}()
}

func (g *GameServer) getBindAddr() string {
	if g.networkServer {
		return "0.0.0.0"
	}
	return "127.0.0.1"
}
