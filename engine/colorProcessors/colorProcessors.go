package colorProcessors

import (
	"github.com/ukoelguanche/graphicsengine/core"
)

type PaletteSwapColorProcessor struct {
	SourcePalette *core.Palette
	TargetPalette *core.Palette
	Sprite        *core.Sprite
}

func (p *PaletteSwapColorProcessor) AppliesTo(sprite *core.Sprite) bool {
	return sprite == p.Sprite
}

func (p *PaletteSwapColorProcessor) ProcessColor(color []byte) []byte {
	i := p.Sprite.CurrentSwapPaletteIndex()
	return p.SourcePalette.ReplacePalette(color, p.TargetPalette, i)
}

type OpacityColorProcessor struct {
	SourcePalette *core.Palette
	TargetPalette *core.Palette
	Sprite        *core.Sprite
}

func (p *OpacityColorProcessor) ProcessColor(color []byte) []byte {
	color[0] = byte(int32(float32(color[1]) * 0.9))
	color[1] = byte(int32(float32(color[1]) * 0.9))
	color[2] = byte(int32(float32(color[1]) * 0.9))
	return color
}
