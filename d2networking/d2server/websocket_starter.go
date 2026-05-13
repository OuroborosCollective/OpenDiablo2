package d2server

import (
	"fmt"
	"net/http"
)

// StartWebSocket starts the WebSocket server on the given port
func (g *GameServer) StartWebSocket(wsPort int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", g.HandleWebSocket)

	g.Infof("Starting WebSocket Game Server @ 0.0.0.0:%d/ws\n", wsPort)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", wsPort), mux)
		if err != nil {
			g.Errorf("WebSocket server failed: %v", err)
		}
	}()
}
