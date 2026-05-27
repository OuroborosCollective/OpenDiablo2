package d2mpq

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/JoshVarga/blast"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2data/d2compression"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2math"
)

// Stream represents a stream of data in an MPQ archive
type Stream struct {
	Data      []byte
	Positions []uint32
	MPQ       *MPQ
	Block     *Block
	Index     uint32
	Size      uint32
	Position  uint32
}

// CreateStream creates an MPQ stream
func CreateStream(mpq *MPQ, block *Block, fileName string) (*Stream, error) {
	s := &Stream{
		MPQ:   mpq,
		Block: block,
		Index: 0xFFFFFFFF, //nolint:gomnd // MPQ magic
	}

	if s.Block.HasFlag(FileFixKey) {
		s.Block.calculateEncryptionSeed(s.MPQ.crypto, fileName)
	}

	s.Size = 0x200 << s.MPQ.header.BlockSize //nolint:gomnd // MPQ magic
	if s.Size == 0 {
		return nil, errors.New("invalid MPQ block size")
	}

	if s.Block.HasFlag(FilePatchFile) {
		return nil, errors.New("patching is not supported")
	}

	if (s.Block.HasFlag(FileCompress) || s.Block.HasFlag(FileImplode)) && !s.Block.HasFlag(FileSingleUnit) {
		if err := s.loadBlockOffsets(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (v *Stream) loadBlockOffsets() error {
	if _, err := v.MPQ.file.Seek(int64(v.Block.FilePosition), io.SeekStart); err != nil {
		return err
	}

	blockPositionCount := ((v.Block.UncompressedFileSize + v.Size - 1) / v.Size) + 1
	if blockPositionCount > 1000000 { // Safety limit for number of blocks
		return errors.New("too many blocks in MPQ file")
	}
	v.Positions = make([]uint32, blockPositionCount)

	if err := binary.Read(v.MPQ.file, binary.LittleEndian, &v.Positions); err != nil {
		return err
	}

	if v.Block.HasFlag(FileEncrypted) {
		v.MPQ.crypto.decrypt(v.Positions, v.Block.EncryptionSeed-1)

		blockPosSize := blockPositionCount << 2 //nolint:gomnd // MPQ magic
		if v.Positions[0] != blockPosSize {
			return errors.New("decryption of MPQ failed")
		}

		if v.Positions[1] > v.Size+blockPosSize {
			return errors.New("decryption of MPQ failed")
		}
	}

	return nil
}

func (v *Stream) Read(buffer []byte, offset, count uint32) (readTotal uint32, err error) {
	if count == 0 {
		return 0, nil
	}

	// Safety check for buffer bounds
	if offset >= uint32(len(buffer)) {
		return 0, io.ErrShortBuffer
	}

	// Prevent reading beyond file size
	if v.Position >= v.Block.UncompressedFileSize {
		return 0, io.EOF
	}

	if v.Block.HasFlag(FileSingleUnit) {
		return v.readInternalSingleUnit(buffer, offset, count)
	}

	var read uint32

	toRead := count
	for toRead > 0 {
		if read, err = v.readInternal(buffer, offset, toRead); err != nil {
			if err == io.EOF && readTotal > 0 {
				return readTotal, nil
			}
			return readTotal, err
		}

		if read == 0 {
			break
		}

		readTotal += read
		offset += read
		toRead -= read

		// If we've reached the end of the file, stop reading
		if v.Position >= v.Block.UncompressedFileSize {
			break
		}
	}

	return readTotal, nil
}

func (v *Stream) readInternalSingleUnit(buffer []byte, offset, count uint32) (uint32, error) {
	if len(v.Data) == 0 {
		if err := v.loadSingleUnit(); err != nil {
			return 0, err
		}
	}

	return v.copy(buffer, offset, v.Position, count)
}

func (v *Stream) readInternal(buffer []byte, offset, count uint32) (uint32, error) {
	if err := v.bufferData(); err != nil {
		return 0, err
	}

	localPosition := v.Position % v.Size

	return v.copy(buffer, offset, localPosition, count)
}

func (v *Stream) copy(buffer []byte, offset, pos, count uint32) (uint32, error) {
	if offset >= uint32(len(buffer)) {
		return 0, io.ErrShortBuffer
	}

	if pos >= uint32(len(v.Data)) {
		return 0, io.EOF
	}

	bytesToCopy := d2math.Min(uint32(len(v.Data))-pos, count)
	if bytesToCopy <= 0 {
		return 0, io.EOF
	}

	if offset+bytesToCopy > uint32(len(buffer)) {
		bytesToCopy = uint32(len(buffer)) - offset
	}

	// Double check bounds to prevent panic
	if offset+bytesToCopy > uint32(len(buffer)) || pos+bytesToCopy > uint32(len(v.Data)) {
		return 0, errors.New("copy bounds exceeded")
	}

	copy(buffer[offset:offset+bytesToCopy], v.Data[pos:pos+bytesToCopy])
	v.Position += bytesToCopy

	return bytesToCopy, nil
}

func (v *Stream) bufferData() (err error) {
	blockIndex := v.Position / v.Size

	if blockIndex == v.Index {
		return nil
	}

	expectedLength := d2math.Min(v.Block.UncompressedFileSize-(blockIndex*v.Size), v.Size)
	if v.Data, err = v.loadBlock(blockIndex, expectedLength); err != nil {
		return err
	}

	v.Index = blockIndex

	return nil
}

func (v *Stream) loadSingleUnit() (err error) {
	if _, err = v.MPQ.file.Seek(int64(v.Block.FilePosition), io.SeekStart); err != nil {
		return err
	}

	fileData := make([]byte, v.Block.CompressedFileSize)

	if _, err = io.ReadFull(v.MPQ.file, fileData); err != nil {
		return err
	}

	if v.Block.HasFlag(FileEncrypted) && v.Block.UncompressedFileSize > 3 {
		if v.Block.EncryptionSeed == 0 {
			return errors.New("unable to determine encryption key")
		}

		v.MPQ.crypto.decryptBytes(fileData, v.Block.EncryptionSeed)
	}

	if v.Block.CompressedFileSize == v.Block.UncompressedFileSize {
		v.Data = fileData
		return nil
	}

	if v.Block.HasFlag(FileCompress) {
		v.Data, err = decompressMulti(fileData, v.Block.UncompressedFileSize)
	} else if v.Block.HasFlag(FileImplode) {
		v.Data, err = pkDecompress(fileData)
	} else {
		v.Data = fileData
	}

	return err
}

func (v *Stream) loadBlock(blockIndex, expectedLength uint32) ([]byte, error) {
	var (
		offset uint32
		toRead uint32
	)

	if v.Block.HasFlag(FileCompress) || v.Block.HasFlag(FileImplode) {
		if blockIndex+1 >= uint32(len(v.Positions)) {
			return []byte{}, errors.New("block index out of bounds")
		}
		offset = v.Positions[blockIndex]

		if v.Positions[blockIndex+1] < offset {
			return []byte{}, errors.New("invalid block offsets in MPQ")
		}

		toRead = v.Positions[blockIndex+1] - offset
	} else {
		offset = blockIndex * v.Size
		toRead = expectedLength
	}

	// Check for potential overflow or absurdly large toRead
	if toRead > 100*1024*1024 { // 100MB limit per block as a safety measure
		return []byte{}, fmt.Errorf("block size too large: %d", toRead)
	}

	offset += v.Block.FilePosition

	// Physical file size boundary check
	fi, err := v.MPQ.file.Stat()
	if err != nil {
		return []byte{}, err
	}
	if int64(offset+toRead) > fi.Size() {
		return []byte{}, fmt.Errorf("read out of bounds: offset %d, size %d, file size %d", offset, toRead, fi.Size())
	}

	data := make([]byte, toRead)

	if _, err := v.MPQ.file.Seek(int64(offset), io.SeekStart); err != nil {
		return []byte{}, err
	}

	if _, err := v.MPQ.file.Read(data); err != nil {
		return []byte{}, err
	}

	if v.Block.HasFlag(FileEncrypted) && v.Block.UncompressedFileSize > 3 {
		if v.Block.EncryptionSeed == 0 {
			return []byte{}, errors.New("unable to determine encryption key")
		}

		v.MPQ.crypto.decryptBytes(data, blockIndex+v.Block.EncryptionSeed)
	}

	if v.Block.HasFlag(FileCompress) {
		if toRead != expectedLength {
			if !v.Block.HasFlag(FileSingleUnit) {
				return decompressMulti(data, expectedLength)
			}

			return pkDecompress(data)
		}

		return data, nil
	}

	if v.Block.HasFlag(FileImplode) {
		if toRead != expectedLength {
			return pkDecompress(data)
		}

		return data, nil
	}

	return data, nil
}

func decompressMulti(data []byte, expectedLength uint32) (res []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("decompression panic: %v", r)
		}
	}()

	if len(data) == 0 {
		return []byte{}, errors.New("empty data for decompression")
	}

	compressionMask := data[0]
	remainingData := data[1:]

	res, err = decompressByMask(compressionMask, remainingData, expectedLength)
	if err != nil {
		return nil, err
	}

	if uint32(len(res)) != expectedLength {
		return nil, fmt.Errorf("decompressed size mismatch: got %d, expected %d", len(res), expectedLength)
	}

	return res, nil
}

func decompressByMask(mask byte, data []byte, expectedLength uint32) ([]byte, error) {
	// Check for combined compression types
	// These are processed in a specific order if multiple bits are set
	// For standard D2 MPQs, these are usually single bits or specific combos like 0x41, 0x81

	// If it's a known combo, handle it
	switch mask {
	case 0x41: // Huffman + Wav Mono
		huff, err := decompressHuffman(data)
		if err != nil {
			return nil, err
		}
		return d2compression.WavDecompress(huff, 1)
	case 0x81: // Huffman + Wav Stereo
		huff, err := decompressHuffman(data)
		if err != nil {
			return nil, err
		}
		return d2compression.WavDecompress(huff, 2)
	}

	// Single algorithms
	switch mask {
	case 1: // Huffman
		return decompressHuffman(data)
	case 2: // ZLib/Deflate
		return deflate(data)
	case 8: // PKLib/Implode
		return pkDecompress(data)
	case 0x10: // BZip2
		return nil, errors.New("bzip2 decompression (0x10) not supported")
	case 0x40: // IMA ADPCM Mono
		return d2compression.WavDecompress(data, 1)
	case 0x80: // IMA ADPCM Stereo
		return d2compression.WavDecompress(data, 2)
	case 0x12:
		return nil, errors.New("lzma decompression (0x12) not supported")
	case 0x22:
		return nil, errors.New("sparse decompression + deflate decompression (0x22) not supported")
	case 0x30:
		return nil, errors.New("sparse decompression + bzip2 decompression (0x30) not supported")
	case 0x48:
		return []byte{}, errors.New("pk + mpqwav decompression (0x48) not supported")
	case 0x88:
		return []byte{}, errors.New("pk + wav decompression (0x88) not supported")
	default:
		return []byte{}, fmt.Errorf("decompression not supported for mask %X", mask)
	}
}

func decompressHuffman(data []byte) ([]byte, error) {
	res := d2compression.HuffmanDecompress(data)
	if res == nil {
		return nil, errors.New("huffman decompression failed")
	}
	return res, nil
}

func deflate(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)

	r, err := zlib.NewReader(b)
	if err != nil {
		return []byte{}, err
	}

	buffer := new(bytes.Buffer)

	_, err = buffer.ReadFrom(r)
	if err != nil {
		return []byte{}, err
	}

	err = r.Close()
	if err != nil {
		return []byte{}, err
	}

	return buffer.Bytes(), nil
}

func pkDecompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)

	r, err := blast.NewReader(b)
	if err != nil {
		return []byte{}, err
	}

	buffer := new(bytes.Buffer)

	if _, err = buffer.ReadFrom(r); err != nil {
		return []byte{}, err
	}

	err = r.Close()
	if err != nil {
		return []byte{}, err
	}

	return buffer.Bytes(), nil
}
