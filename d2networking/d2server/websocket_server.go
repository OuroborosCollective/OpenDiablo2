package d2server

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2netpacket"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2netpacket/d2netpackettype"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2server/d2wsclientconnection"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for the boilerplate
	},
}

func (g *GameServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		g.Errorf("failed to upgrade to websocket: %v", err)
		return
	}

	g.Infof("Accepting WebSocket connection: %s\n", conn.RemoteAddr().String())

	defer func() {
		if err := conn.Close(); err != nil {
			g.Errorf("failed to close the websocket connection: %s\n", conn.RemoteAddr())
		}
	}()

	var client *d2wsclientconnection.WSClientConnection
	connected := false

	for {
		var packet d2netpacket.NetPacket
		err := conn.ReadJSON(&packet)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				g.Errorf("websocket error: %v", err)
			}
			if client != nil {
				g.OnClientDisconnected(client)
			}
			break
		}

		if !connected {
			if packet.PacketType != d2netpackettype.PlayerConnectionRequest {
				g.Infof("Closing connection with %s: did not receive new player connection request...", conn.RemoteAddr().String())
				return
			}

			// We need a modified registerConnection or similar for WS
			client, err = g.registerWSConnection(packet.PacketData, conn)
			if err != nil {
				return
			}

			connected = true
		}

		select {
		case <-g.ctx.Done():
			return
		default:
			g.packetManagerChan <- ReceivedPacket{
				Client: client,
				Packet: packet,
			}
		}
	}
}

func (g *GameServer) registerWSConnection(b []byte, conn *websocket.Conn) (*d2wsclientconnection.WSClientConnection, error) {
	g.Lock()
	defer g.Unlock()

	// check to see if the server is full
	if len(g.connections) >= g.maxConnections {
		// Send server full packet (simplified for WS)
		sf, _ := d2netpacket.CreateServerFullPacket()
		conn.WriteJSON(sf)
		return nil, errServerFull
	}

	packet, err := d2netpacket.UnmarshalPlayerConnectionRequest(b)
	if err != nil {
		g.Errorf("Failed to unmarshal PlayerConnectionRequest: %s\n", err)
		return nil, err
	}

	if _, ok := g.connections[packet.ID]; ok {
		g.Errorf("%v", errPlayerAlreadyExists)
		return nil, errPlayerAlreadyExists
	}

	client := d2wsclientconnection.CreateWSClientConnection(conn, packet.ID)
	client.SetPlayerState(packet.PlayerState)

	// Since OnClientConnected takes the interface, this works
	g.OnClientConnected(client)

	return client, nil
}
