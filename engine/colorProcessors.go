package engine

import (
	"time"

	"github.com/ukoelguanche/graphicsengine/core"
	"github.com/ukoelguanche/graphicsengine/drivers"
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

type FadeToBlack struct {
	drivers.PixelTransformer
	StartTime time.Time
	Duration  time.Duration
}

func (ftb *FadeToBlack) Transform(pixels []byte) {
	elapsed := time.Since(ftb.StartTime)
	duration := ftb.Duration

	t := elapsed.Seconds() / duration.Seconds()

	if t > 1.0 {
		t = 1.0
	}

	easedT := EaseInOutCubic(t)

	start := 1.0
	target := 0.0

	nextY := float32(start + (target-start)*easedT)

	for i := 0; i < len(pixels); i += 4 {
		pixels[i] = uint8(float32(pixels[i]) * nextY)
		pixels[i+1] = uint8(float32(pixels[i+1]) * nextY)
		pixels[i+2] = uint8(float32(pixels[i+2]) * nextY)
	}
}

type PixelFade struct {
	drivers.PixelTransformer
	StartTime time.Time
	Duration  time.Duration
	GridSize  int
}

func (pf *PixelFade) Transform(pixels []byte) {
	elapsed := time.Since(pf.StartTime)
	duration := pf.Duration

	t := elapsed.Seconds() / duration.Seconds()

	if t > 1.0 {
		t = 1.0
	}

	start := 0.0
	target := float64(pf.GridSize) * 2

	step := int(float32(start + (target-start)*t))

	for i := 0; i < drivers.VH; i++ { // 0 to 200: rows
		for j := 0; j < drivers.VW; j++ { // 0 to 320*4: columns (4 bytes per pixel)
			ii := i % pf.GridSize
			jj := j % pf.GridSize

			if jj <= step-ii {
				pixelPos := (i*drivers.VW + j) * 4
				pixels[pixelPos] = 0
				pixels[pixelPos+1] = 0
				pixels[pixelPos+2] = 0
			}
		}
	}
}
