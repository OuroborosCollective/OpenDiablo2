package ebiten

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

func d2ToEbitenFilter(filter d2enum.Filter) ebiten.Filter {
	switch filter {
	case d2enum.FilterNearest:
		return ebiten.FilterNearest
	default:
		return ebiten.FilterLinear
	}
}
