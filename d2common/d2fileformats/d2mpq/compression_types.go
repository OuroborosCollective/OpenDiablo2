package d2mpq

// Compression types for MPQ files
const (
	CompHuffman      uint8 = 0x01
	CompDeflate      uint8 = 0x02
	CompPKLib        uint8 = 0x08
	CompBZip2        uint8 = 0x10
	CompSparse       uint8 = 0x20
	CompADPCMMono    uint8 = 0x40
	CompADPCMStereo  uint8 = 0x80
	CompLZMA         uint8 = 0x12
)
