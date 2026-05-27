package d2mpq

import (
	"errors"
	"io"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

var _ d2interface.DataStream = &MpqDataStream{} // Static check to confirm struct conforms to interface

// MpqDataStream represents a stream for MPQ data.
type MpqDataStream struct {
	stream *Stream
}

// Read reads data from the data stream
func (m *MpqDataStream) Read(p []byte) (n int, err error) {
	totalRead, err := m.stream.Read(p, 0, uint32(len(p)))
	return int(totalRead), err
}

// Seek sets the position of the data stream
func (m *MpqDataStream) Seek(offset int64, whence int) (int64, error) {
	var newPos int64

	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekCurrent:
		newPos = int64(m.stream.Position) + offset
	case io.SeekEnd:
		newPos = int64(m.stream.Block.UncompressedFileSize) + offset
	default:
		return 0, errors.New("invalid whence")
	}

	if newPos < 0 {
		return 0, errors.New("negative position")
	}

	m.stream.Position = uint32(newPos)

	return int64(m.stream.Position), nil
}

// Close closes the data stream
func (m *MpqDataStream) Close() error {
	m.stream = nil
	return nil
}
