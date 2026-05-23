package d2netpacket

import (
	"encoding/json"

	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2netpacket/d2netpackettype"
)

// AssetMetadata represents metadata for an asset.
type AssetMetadata struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Size string `json:"size"`
	Path string `json:"path"`
}

// AssetListPacket contains a list of loaded assets.
type AssetListPacket struct {
	Assets []AssetMetadata `json:"assets"`
}

// CreateAssetListPacket returns a NetPacket which declares an
// AssetListPacket with the given assets.
func CreateAssetListPacket(assets []AssetMetadata) (NetPacket, error) {
	packet := AssetListPacket{
		Assets: assets,
	}

	b, err := json.Marshal(packet)
	if err != nil {
		return NetPacket{PacketType: d2netpackettype.AssetList}, err
	}

	return NetPacket{
		PacketType: d2netpackettype.AssetList,
		PacketData: b,
	}, nil
}

// UnmarshalAssetList unmarshals the data to an AssetListPacket struct
func UnmarshalAssetList(packet []byte) (AssetListPacket, error) {
	var resp AssetListPacket

	if err := json.Unmarshal(packet, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
