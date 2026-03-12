package engine

import "github.com/ukoelguanche/graphicsengine/core"

type PaletteSwapColorProcessor struct {
	SourcePalette *core.Palette
	TargetPalette *core.Palette
	Sprite        *core.Sprite
}

func (p *PaletteSwapColorProcessor) ProcessColor(color []byte) []byte {
	i := p.Sprite.CurrentSwapPaletteIndex()
	return p.SourcePalette.ReplacePalette(color, p.TargetPalette, i)
}
