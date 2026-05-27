package d2dcc

import (
	"testing"
)

func TestLoad_InvalidSignature(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00, 0x00}
	_, err := Load(data)
	if err == nil {
		t.Error("Expected error for invalid signature, got nil")
	}
	expectedErr := "signature expected to be 0x74 but it is not"
	if err.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestLoad_InvalidHeaderValue(t *testing.T) {
	// Sig: 0x74
	// Version: 1
	// Directions: 1
	// Frames: 1
	// MustBeOne: 0 (invalid)
	data := []byte{0x74, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	_, err := Load(data)
	if err == nil {
		t.Error("Expected error for invalid MustBeOne value, got nil")
	}
	expectedErr := "this value isn't 1. It has to be 1"
	if err.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestClone(t *testing.T) {
	dcc := &DCC{
		Signature:          0x74,
		Version:            1,
		NumberOfDirections: 1,
		FramesPerDirection: 1,
		Directions: []*DCCDirection{
			{OutSizeCoded: 100},
		},
		directionOffsets: []int{10},
		fileData:         []byte{0xDE, 0xAD, 0xBE, 0xEF},
	}

	clone := dcc.Clone()

	// Basic checks
	if clone.Signature != dcc.Signature {
		t.Errorf("Signature mismatch: %v != %v", clone.Signature, dcc.Signature)
	}

	// check directionOffsets independence
	dcc.directionOffsets[0] = 20
	if clone.directionOffsets[0] == 20 {
		t.Error("directionOffsets should be deep copied, but changes to original affect clone")
	}

	// check fileData independence
	dcc.fileData[0] = 0x00
	if clone.fileData[0] == 0x00 {
		t.Error("fileData should be deep copied, but changes to original affect clone")
	}

	// check Directions independence
	if len(clone.Directions) != len(dcc.Directions) {
		t.Fatalf("Directions length mismatch: %v != %v", len(clone.Directions), len(dcc.Directions))
	}
	if clone.Directions[0] == nil {
		t.Fatal("clone.Directions[0] is nil")
	}
	if clone.Directions[0] == dcc.Directions[0] {
		t.Error("Directions elements should be deep copied, but they point to the same memory")
	}

	dcc.Directions[0].OutSizeCoded = 200
	if clone.Directions[0].OutSizeCoded == 200 {
		t.Error("Directions content should be deep copied, but changes to original affect clone")
	}
}
