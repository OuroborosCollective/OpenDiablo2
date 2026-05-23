package d2netpacket

import (
	"encoding/json"

	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2netpacket/d2netpackettype"
)

// AxiomaticStatusPacket contains the current state of the Axiomatic engine.
type AxiomaticStatusPacket struct {
	Resonance float64 `json:"resonance"`
	Cycle     float64 `json:"cycle"`
}

// CreateAxiomaticStatusPacket returns a NetPacket which declares an
// AxiomaticStatusPacket with the given resonance and cycle.
func CreateAxiomaticStatusPacket(resonance, cycle float64) (NetPacket, error) {
	status := AxiomaticStatusPacket{
		Resonance: resonance,
		Cycle:     cycle,
	}

	b, err := json.Marshal(status)
	if err != nil {
		return NetPacket{PacketType: d2netpackettype.AxiomaticStatus}, err
	}

	return NetPacket{
		PacketType: d2netpackettype.AxiomaticStatus,
		PacketData: b,
	}, nil
}

// UnmarshalAxiomaticStatus unmarshals the data to an AxiomaticStatusPacket struct
func UnmarshalAxiomaticStatus(packet []byte) (AxiomaticStatusPacket, error) {
	var resp AxiomaticStatusPacket

	if err := json.Unmarshal(packet, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
