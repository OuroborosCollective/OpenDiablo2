package d2dcc

import (
	"testing"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2geom"
)

func TestDCCDirection_calculateCells(t *testing.T) {
	tests := []struct {
		name          string
		width, height int
		wantH, wantV  int
	}{
		{"4x4", 4, 4, 1, 1},
		{"5x4", 5, 4, 2, 1},
		{"4x5", 4, 5, 1, 2},
		{"8x8", 8, 8, 2, 2},
		{"1x1", 1, 1, 1, 1},
		{"2x2", 2, 2, 1, 1},
		{"9x9", 9, 9, 3, 3},
		{"12x12", 12, 12, 3, 3},
		{"13x13", 13, 13, 4, 4},
		{"0x0", 0, 0, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &DCCDirection{
				Box: d2geom.Rectangle{Width: tt.width, Height: tt.height},
			}
			v.calculateCells()
			if v.HorizontalCellCount != tt.wantH {
				t.Errorf("HorizontalCellCount = %v, want %v", v.HorizontalCellCount, tt.wantH)
			}
			if v.VerticalCellCount != tt.wantV {
				t.Errorf("VerticalCellCount = %v, want %v", v.VerticalCellCount, tt.wantV)
			}
			if len(v.Cells) != tt.wantH*tt.wantV {
				t.Errorf("len(Cells) = %v, want %v", len(v.Cells), tt.wantH*tt.wantV)
			}

			// Verify cell dimensions
			for i, cell := range v.Cells {
				if cell == nil {
					t.Fatalf("Cell[%d] is nil", i)
				}

				cellX := i % v.HorizontalCellCount
				cellY := i / v.HorizontalCellCount

				expectedWidth := 4
				if cellX == v.HorizontalCellCount-1 {
					if v.HorizontalCellCount == 1 {
						expectedWidth = tt.width
					} else {
						expectedWidth = tt.width - (4 * (v.HorizontalCellCount - 1))
					}
				}

				expectedHeight := 4
				if cellY == v.VerticalCellCount-1 {
					if v.VerticalCellCount == 1 {
						expectedHeight = tt.height
					} else {
						expectedHeight = tt.height - (4 * (v.VerticalCellCount - 1))
					}
				}

				if cell.Width != expectedWidth {
					t.Errorf("Cell[%d].Width = %d, want %d", i, cell.Width, expectedWidth)
				}
				if cell.Height != expectedHeight {
					t.Errorf("Cell[%d].Height = %d, want %d", i, cell.Height, expectedHeight)
				}

				if cell.XOffset != cellX*4 {
					t.Errorf("Cell[%d].XOffset = %d, want %d", i, cell.XOffset, cellX*4)
				}
				if cell.YOffset != cellY*4 {
					t.Errorf("Cell[%d].YOffset = %d, want %d", i, cell.YOffset, cellY*4)
				}
			}
		})
	}
}

func TestDCCDirection_verify(t *testing.T) {
	// Create bitmunchers with some bits read
	data := make([]byte, 16)

	v := &DCCDirection{
		EqualCellsBitstreamSize:    1,
		PixelMaskBitstreamSize:     2,
		EncodingTypeBitsreamSize:   3,
		RawPixelCodesBitstreamSize: 4,
	}

	ec := d2datautils.CreateBitMuncher(data, 0)
	ec.GetBit()

	pm := d2datautils.CreateBitMuncher(data, 0)
	pm.GetBits(2)

	et := d2datautils.CreateBitMuncher(data, 0)
	et.GetBits(3)

	rp := d2datautils.CreateBitMuncher(data, 0)
	rp.GetBits(4)

	// Should not panic
	v.verify(ec, pm, et, rp)
}

func TestDCCDirection_verify_Panic(t *testing.T) {
	data := make([]byte, 16)

	v := &DCCDirection{
		EqualCellsBitstreamSize:    1,
		PixelMaskBitstreamSize:     2,
		EncodingTypeBitsreamSize:   3,
		RawPixelCodesBitstreamSize: 4,
	}

	tests := []struct {
		name string
		ec   int
		pm   int
		et   int
		rp   int
	}{
		{"wrong ec", 0, 2, 3, 4},
		{"wrong pm", 1, 0, 3, 4},
		{"wrong et", 1, 2, 0, 4},
		{"wrong rp", 1, 2, 3, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("verify did not panic")
				}
			}()

			ec := d2datautils.CreateBitMuncher(data, 0)
			ec.GetBits(tt.ec)

			pm := d2datautils.CreateBitMuncher(data, 0)
			pm.GetBits(tt.pm)

			et := d2datautils.CreateBitMuncher(data, 0)
			et.GetBits(tt.et)

			rp := d2datautils.CreateBitMuncher(data, 0)
			rp.GetBits(tt.rp)

			v.verify(ec, pm, et, rp)
		})
	}
}
