package engine

import (
	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/core"
	"github.com/ukoelguanche/graphicsengine/drivers"
)

type CachedSpriteDraw struct {
	Sprite   *core.Sprite
	Frame    core.Frame
	Position core.Point
}

type CachedLayer struct {
	BaseMovable
	sprite   *core.Sprite
	isStatic bool
	easeFunc func(float64) float64
}

func (l *CachedLayer) GetSprite() *core.Sprite                 { return l.sprite }
func (l *CachedLayer) GetEaseFunction() func(float64) float64  { return l.easeFunc }
func (l *CachedLayer) SetEaseFunction(f func(float64) float64) { l.easeFunc = f }

func (l *CachedLayer) Draw(d interfaces.Drawer) {
	d.DrawSpriteRect(l.sprite, l.sprite.Frames[0], l.position)
}

func (l *CachedLayer) Update()        { l.UpdatePosition(l) }
func (l *CachedLayer) IsStatic() bool { return l.isStatic }

func BuildCachedLayer(name string, draws []CachedSpriteDraw, isStatic bool) *CachedLayer {
	bitmap := buildCachedBitmap(name, draws)

	for _, draw := range draws {
		blitFrame(bitmap, draw.Sprite, draw.Frame, draw.Position)
	}

	sprite := &core.Sprite{
		Name:      name,
		Bitmap:    bitmap,
		Frames:    []core.Frame{{Point: core.Point{X: 0, Y: 0}, Size: core.Size{W: drivers.VW, H: drivers.VH}}},
		Sequences: map[string]core.Sequence{"idle": {0}},
	}

	return &CachedLayer{
		BaseMovable: BaseMovable{
			position:       core.Point{X: 0, Y: 0},
			startPosition:  core.Point{X: 0, Y: 0},
			targetPosition: core.Point{X: 0, Y: 0},
			moving:         false,
		},
		sprite:   sprite,
		isStatic: isStatic,
	}
}

func buildCachedBitmap(name string, draws []CachedSpriteDraw) *core.Bitmap {
	bitmap := &core.Bitmap{
		Name: name,
		W:    drivers.VW,
		H:    drivers.VH,
	}

	for _, draw := range draws {
		source := draw.Sprite.GetBitmap()
		if len(source.IndexedPixels) > 0 && source.Palette != nil {
			bitmap.IndexedPixels = make([]byte, drivers.VW*drivers.VH)
			bitmap.Palette = source.Palette
			return bitmap
		}
	}

	bitmap.Pixels = make([]byte, drivers.VW*drivers.VH*4)
	return bitmap
}

func blitFrame(target *core.Bitmap, sprite *core.Sprite, frame core.Frame, position core.Point) {
	source := sprite.GetBitmap()

	for sy := 0; sy < int(frame.Size.H); sy++ {
		dstY := int(position.Y) + sy
		if dstY < 0 || dstY >= int(target.H) {
			continue
		}

		for sx := 0; sx < int(frame.Size.W); sx++ {
			dstX := int(position.X) + sx
			if dstX < 0 || dstX >= int(target.W) {
				continue
			}

			srcX := int(frame.Point.X) + sx
			srcY := int(frame.Point.Y) + sy
			if srcX < 0 || srcX >= int(source.W) || srcY < 0 || srcY >= int(source.H) {
				continue
			}

			srcPixelIndex := srcY*int(source.W) + srcX
			dstPixelIndex := dstY*int(target.W) + dstX
			blitPixel(target, source, srcPixelIndex, dstPixelIndex)
		}
	}
}

func blitPixel(target *core.Bitmap, source *core.Bitmap, srcPixelIndex int, dstPixelIndex int) {
	if len(target.IndexedPixels) > 0 && len(source.IndexedPixels) > 0 && target.Palette == source.Palette {
		paletteIndex := source.IndexedPixels[srcPixelIndex]
		if paletteIndex == 0 {
			return
		}

		target.IndexedPixels[dstPixelIndex] = paletteIndex
		return
	}

	srcOffset := srcPixelIndex * 4
	if source.Pixels[srcOffset+3] < 128 {
		return
	}

	dstOffset := dstPixelIndex * 4
	copy(target.Pixels[dstOffset:dstOffset+4], source.Pixels[srcOffset:srcOffset+4])
}
