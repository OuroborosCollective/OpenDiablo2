package d2netpacket

import (
	"encoding/json"

	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2netpacket/d2netpackettype"
)

// AssetMetadata represents metadata for a single asset.
type AssetMetadata struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Size string `json:"size"`
	Path string `json:"path"`
}

// AssetMetadataListPacket contains a list of AssetMetadata.
type AssetMetadataListPacket struct {
	Assets []AssetMetadata `json:"assets"`
}

// CreateAssetMetadataListPacket returns a NetPacket declaring an AssetMetadataListPacket.
func CreateAssetMetadataListPacket(assets []AssetMetadata) (NetPacket, error) {
	packetData := AssetMetadataListPacket{
		Assets: assets,
	}

	b, err := json.Marshal(packetData)
	if err != nil {
		return NetPacket{PacketType: d2netpackettype.AssetMetadataList}, err
	}

	return NetPacket{
		PacketType: d2netpackettype.AssetMetadataList,
		PacketData: b,
	}, nil
}

// UnmarshalAssetMetadataList unmarshals data to an AssetMetadataListPacket struct.
func UnmarshalAssetMetadataList(packet []byte) (AssetMetadataListPacket, error) {
	var p AssetMetadataListPacket
	if err := json.Unmarshal(packet, &p); err != nil {
		return p, err
	}

	return p, nil
}
