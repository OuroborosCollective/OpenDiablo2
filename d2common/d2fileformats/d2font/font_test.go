package d2font

import (
	"image/color"
	"testing"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font/d2fontglyph"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

type mockAnimation struct {
	d2interface.Animation
	colorMod     color.Color
	currentFrame int
	renderCalled int
}

func (m *mockAnimation) SetColorMod(c color.Color) {
	m.colorMod = c
}

func (m *mockAnimation) SetCurrentFrame(idx int) error {
	m.currentFrame = idx
	return nil
}

func (m *mockAnimation) Render(target d2interface.Surface) {
	m.renderCalled++
}

func (m *mockAnimation) GetFrameBounds() (int, int) {
	return 10, 20
}

type mockSurface struct {
	d2interface.Surface
	translations [][2]int
	pops         int
}

func (m *mockSurface) PushTranslation(x, y int) {
	m.translations = append(m.translations, [2]int{x, y})
}

func (m *mockSurface) PopN(n int) {
	m.pops += n
}

func TestLoad(t *testing.T) {
	// Load skips 5 (sig) + 8 (unknown) = 13 bytes.
	header := []byte("Woo!\x01")
	header = append(header, make([]byte, 8)...)

	glyph1 := []byte{0x41, 0x00} // 'A'
	glyph1 = append(glyph1, 0)
	glyph1 = append(glyph1, 8)
	glyph1 = append(glyph1, 12)
	glyph1 = append(glyph1, []byte{1, 0, 0}...)
	glyph1 = append(glyph1, []byte{0, 0}...)
	glyph1 = append(glyph1, []byte{0, 0, 0, 0}...)

	data := append(header, glyph1...)

	font, err := Load(data)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(font.Glyphs) != 1 {
		t.Errorf("Expected 1 glyph, got %d", len(font.Glyphs))
	}

	g, ok := font.Glyphs['A']
	if !ok {
		t.Fatal("Glyph 'A' not found")
	}

	if g.Width() != 8 || g.Height() != 12 {
		t.Errorf("Unexpected glyph size: %dx%d", g.Width(), g.Height())
	}
}

func TestLoad_Truncated(t *testing.T) {
	data := []byte("Woo!")
	_, err := Load(data)
	if err == nil {
		t.Error("Expected error for truncated signature, got nil")
	}
}

func TestLoad_InvalidSignature(t *testing.T) {
	data := []byte("BadSig")
	_, err := Load(data)
	if err == nil {
		t.Error("Expected error for invalid signature, got nil")
	}
}

func TestMarshalRoundTrip(t *testing.T) {
	font := &Font{
		color: color.White,
		Glyphs: map[rune]*d2fontglyph.FontGlyph{
			'A': d2fontglyph.Create(0, 8, 12),
		},
	}

	marshaled := font.Marshal()

	font2, err := Load(marshaled)
	if err != nil {
		t.Fatalf("Load of marshaled data failed: %v", err)
	}

	if _, ok := font2.Glyphs['A']; !ok {
		t.Error("Glyph 'A' not found in round-trip")
	}
}

func TestSetBackground(t *testing.T) {
	f := &Font{
		Glyphs: map[rune]*d2fontglyph.FontGlyph{
			'A': d2fontglyph.Create(0, 8, 12),
		},
	}

	anim := &mockAnimation{}
	f.SetBackground(anim)

	if f.sheet != anim {
		t.Error("SetBackground did not set sheet")
	}

	g := f.Glyphs['A']
	if g.Height() != 20 {
		t.Errorf("Expected height 20 from animation, got %d", g.Height())
	}
}

func TestSetColor(t *testing.T) {
	f := &Font{}
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	f.SetColor(c)

	if f.color != c {
		t.Error("SetColor did not set color")
	}
}

func TestGetTextMetrics(t *testing.T) {
	f := &Font{
		Glyphs: map[rune]*d2fontglyph.FontGlyph{
			'A': d2fontglyph.Create(0, 10, 20),
			'B': d2fontglyph.Create(1, 15, 25),
		},
	}

	// Single line
	w, h := f.GetTextMetrics("AB")
	if w != 25 || h != 25 {
		t.Errorf("Metrics for 'AB' failed: expected 25x25, got %dx%d", w, h)
	}

	// Multi line
	w, h = f.GetTextMetrics("A\nB")
	if w != 15 || h != 45 {
		t.Errorf("Metrics for 'A\nB' failed: expected 15x45, got %dx%d", w, h)
	}

	// Empty string
	w, h = f.GetTextMetrics("")
	if w != 0 || h != 0 {
		t.Errorf("Metrics for '' failed: expected 0x0, got %dx%d", w, h)
	}
}

func TestRenderText(t *testing.T) {
	f := &Font{
		color: color.White,
		Glyphs: map[rune]*d2fontglyph.FontGlyph{
			'A': d2fontglyph.Create(5, 10, 20),
		},
	}
	anim := &mockAnimation{}
	f.sheet = anim

	surf := &mockSurface{}

	err := f.RenderText("A\nA", surf)
	if err != nil {
		t.Fatalf("RenderText failed: %v", err)
	}

	if anim.colorMod != color.White {
		t.Errorf("Expected color mod White, got %v", anim.colorMod)
	}

	if anim.renderCalled != 2 {
		t.Errorf("Expected 2 render calls, got %d", anim.renderCalled)
	}

	if anim.currentFrame != 5 {
		t.Errorf("Expected frame 5, got %d", anim.currentFrame)
	}

	expectedTranslations := [][2]int{
		{10, 0},
		{0, 20},
		{10, 0},
		{0, 20},
	}

	if len(surf.translations) != len(expectedTranslations) {
		t.Fatalf("Expected %d translations, got %d", len(expectedTranslations), len(surf.translations))
	}

	for i, tr := range expectedTranslations {
		if surf.translations[i] != tr {
			t.Errorf("Translation %d mismatch: expected %v, got %v", i, tr, surf.translations[i])
		}
	}

	if surf.pops != 4 {
		t.Errorf("Expected 4 pops, got %d", surf.pops)
	}
}
