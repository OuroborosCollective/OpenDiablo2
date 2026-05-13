package d2wsclientconnection

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2hero"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2client/d2clientconnectiontype"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2netpacket"
)

type WSClientConnection struct {
	sync.Mutex
	conn        *websocket.Conn
	id          string
	playerState *d2hero.HeroState
}

func CreateWSClientConnection(conn *websocket.Conn, id string) *WSClientConnection {
	return &WSClientConnection{
		conn: conn,
		id:   id,
	}
}

func (c *WSClientConnection) GetUniqueID() string {
	return c.id
}

func (c *WSClientConnection) GetConnectionType() d2clientconnectiontype.ClientConnectionType {
	return d2clientconnectiontype.LANClient
}

func (c *WSClientConnection) SendPacketToClient(packet d2netpacket.NetPacket) error {
	c.Lock()
	defer c.Unlock()

	return c.conn.WriteJSON(packet)
}

func (c *WSClientConnection) GetPlayerState() *d2hero.HeroState {
	return c.playerState
}

func (c *WSClientConnection) SetPlayerState(playerState *d2hero.HeroState) {
	c.playerState = playerState
}
