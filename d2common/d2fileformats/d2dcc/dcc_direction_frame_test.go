package d2dcc

import (
	"testing"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2geom"
)

func TestCreateDCCDirectionFrame(t *testing.T) {
	// Setup DCCDirection with bit counts
	direction := &DCCDirection{
		Variable0Bits:    1,
		WidthBits:       8,
		HeightBits:      8,
		XOffsetBits:     8,
		YOffsetBits:     8,
		OptionalDataBits: 4,
		CodedBytesBits:   8,
	}

	data := make([]byte, 10)
	bm := d2datautils.CreateBitMuncher(data, 0)
	frame := CreateDCCDirectionFrame(bm, direction)

	if frame == nil {
		t.Fatal("Expected frame to be created")
	}

	if !frame.valid {
		t.Error("Expected frame to be valid")
	}

	if frame.FrameIsBottomUp {
		t.Error("Expected FrameIsBottomUp to be false")
	}
}

func TestCreateDCCDirectionFrame_PanicBottomUp(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("CreateDCCDirectionFrame did not panic on BottomUp frame")
		}
	}()

	direction := &DCCDirection{
		Variable0Bits: 0, WidthBits: 0, HeightBits: 0, XOffsetBits: 0, YOffsetBits: 0, OptionalDataBits: 0, CodedBytesBits: 0,
	}

	// Set the 1st bit to 1 (BottomUp).
	// BitMuncher reads LSB first.
	data := []byte{0x01}
	bm := d2datautils.CreateBitMuncher(data, 0)
	CreateDCCDirectionFrame(bm, direction)
}

func TestRecalculateCells(t *testing.T) {
	tests := []struct {
		name                string
		directionBox        d2geom.Rectangle
		frameBox            d2geom.Rectangle
		wantHorizontalCells int
		wantVerticalCells   int
	}{
		{
			"Aligned 8x8",
			d2geom.Rectangle{Left: 0, Top: 0, Width: 100, Height: 100},
			d2geom.Rectangle{Left: 0, Top: 0, Width: 8, Height: 8},
			2, 2,
		},
		{
			"Offset 1,1 8x8",
			d2geom.Rectangle{Left: 0, Top: 0, Width: 100, Height: 100},
			d2geom.Rectangle{Left: 1, Top: 1, Width: 8, Height: 8},
			3, 3,
		},
		{
			"Small 2x2 aligned",
			d2geom.Rectangle{Left: 0, Top: 0, Width: 100, Height: 100},
			d2geom.Rectangle{Left: 0, Top: 0, Width: 2, Height: 2},
			1, 1,
		},
		{
			"Small 2x2 offset 3,3",
			d2geom.Rectangle{Left: 0, Top: 0, Width: 100, Height: 100},
			d2geom.Rectangle{Left: 3, Top: 3, Width: 2, Height: 2},
			2, 2,
		},
		{
			"Negative offset 8x8 at -4,-4",
			d2geom.Rectangle{Left: 0, Top: 0, Width: 100, Height: 100},
			d2geom.Rectangle{Left: -4, Top: -4, Width: 8, Height: 8},
			2, 2,
		},
		{
			"Negative offset 8x8 at -1,-1",
			d2geom.Rectangle{Left: 0, Top: 0, Width: 100, Height: 100},
			d2geom.Rectangle{Left: -1, Top: -1, Width: 8, Height: 8},
			3, 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			direction := &DCCDirection{Box: tt.directionBox}
			frame := &DCCDirectionFrame{
				Width:  tt.frameBox.Width,
				Height: tt.frameBox.Height,
				Box:    tt.frameBox,
			}

			frame.recalculateCells(direction)

			if frame.HorizontalCellCount != tt.wantHorizontalCells {
				t.Errorf("HorizontalCellCount = %d, want %d", frame.HorizontalCellCount, tt.wantHorizontalCells)
			}
			if frame.VerticalCellCount != tt.wantVerticalCells {
				t.Errorf("VerticalCellCount = %d, want %d", frame.VerticalCellCount, tt.wantVerticalCells)
			}
			if len(frame.Cells) != frame.HorizontalCellCount*frame.VerticalCellCount {
				t.Errorf("len(Cells) = %d, want %d", len(frame.Cells), frame.HorizontalCellCount*frame.VerticalCellCount)
			}
		})
	}
}
