package d2mpq

import (
	"encoding/binary"
	"io"
	"strings"
)

type crypto struct {
	buffer [0x500]uint32
}

func newCrypto() *crypto {
	c := &crypto{}
	c.initialize()

	return c
}

func (c *crypto) lookup(index uint32) uint32 {
	return c.buffer[index]
}

//nolint:gomnd // Decryption magic
func (c *crypto) initialize() {
	seed := uint32(0x00100001)

	for index1 := 0; index1 < 0x100; index1++ {
		index2 := index1

		for i := 0; i < 5; i++ {
			seed = (seed*125 + 3) % 0x2AAAAB
			temp1 := (seed & 0xFFFF) << 0x10
			seed = (seed*125 + 3) % 0x2AAAAB
			temp2 := seed & 0xFFFF
			c.buffer[index2] = temp1 | temp2
			index2 += 0x100
		}
	}
}

//nolint:gomnd // Decryption magic
func (c *crypto) decrypt(data []uint32, seed uint32) {
	seed2 := uint32(0xeeeeeeee)

	for i := 0; i < len(data); i++ {
		seed2 += c.lookup(0x400 + (seed & 0xff))
		result := data[i]
		result ^= seed + seed2

		seed = ((^seed << 21) + 0x11111111) | (seed >> 11)
		seed2 = result + seed2 + (seed2 << 5) + 3
		data[i] = result
	}
}

//nolint:gomnd // Decryption magic
func (c *crypto) decryptBytes(data []byte, seed uint32) {
	seed2 := uint32(0xEEEEEEEE)

	for i := 0; i < len(data)-3; i += 4 {
		seed2 += c.lookup(0x400 + (seed & 0xFF))
		result := binary.LittleEndian.Uint32(data[i : i+4])
		result ^= seed + seed2
		seed = ((^seed << 21) + 0x11111111) | (seed >> 11)
		seed2 = result + seed2 + (seed2 << 5) + 3

		data[i+0] = uint8(result & 0xff)
		data[i+1] = uint8((result >> 8) & 0xff)
		data[i+2] = uint8((result >> 16) & 0xff)
		data[i+3] = uint8((result >> 24) & 0xff)
	}
}

//nolint:gomnd // Decryption magic
func (c *crypto) decryptTable(r io.Reader, size uint32, name string) ([]uint32, error) {
	seed := c.hashString(name, 3)
	seed2 := uint32(0xEEEEEEEE)
	size *= 4

	table := make([]uint32, size)
	buf := make([]byte, 4)

	for i := uint32(0); i < size; i++ {
		seed2 += c.buffer[0x400+(seed&0xff)]

		if _, err := r.Read(buf); err != nil {
			return table, err
		}

		result := binary.LittleEndian.Uint32(buf)
		result ^= seed + seed2

		seed = ((^seed << 21) + 0x11111111) | (seed >> 11)
		seed2 = result + seed2 + (seed2 << 5) + 3
		table[i] = result
	}

	return table, nil
}

func (c *crypto) hashFilename(key string) uint64 {
	a, b := c.hashString(key, 1), c.hashString(key, 2)

	return uint64(a)<<32 | uint64(b)
}

//nolint:gomnd // Decryption magic
func (c *crypto) hashString(key string, hashType uint32) uint32 {
	seed1 := uint32(0x7FED7FED)
	seed2 := uint32(0xEEEEEEEE)

	/* prepare seeds. */
	for _, char := range strings.ToUpper(key) {
		seed1 = c.lookup((hashType*0x100)+uint32(char)) ^ (seed1 + seed2)
		seed2 = uint32(char) + seed1 + seed2 + (seed2 << 5) + 3
	}

	return seed1
}

//nolint:unused,deadcode,gomnd // will use this for creating mpq's
func (c *crypto) encrypt(data []uint32, seed uint32) {
	seed2 := uint32(0xeeeeeeee)

	for i := 0; i < len(data); i++ {
		seed2 += c.lookup(0x400 + (seed & 0xff))
		result := data[i]
		result ^= seed + seed2

		seed = ((^seed << 21) + 0x11111111) | (seed >> 11)
		seed2 = data[i] + seed2 + (seed2 << 5) + 3
		data[i] = result
	}
}
