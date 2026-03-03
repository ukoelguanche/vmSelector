package drivers

import "apodeiktikos.com/fbtest/model"

const (
	VW, VH = 320, 200
	SW, SH = 1280, 720
)

var GlobalDisplay *Display

func DrawString(sprite *model.Sprite, text string, x, y int32, typography string) {
	cursorX := x

	letters := sprite.GetSection(typography)

	for _, char := range text {
		sChar := string(char)
		rect, ok := letters[sChar]
		if !ok {
			cursorX += 8
			continue
		}

		GlobalDisplay.DrawSpriteRect(sprite.Bitmap, rect, cursorX, y)
		cursorX += int32(rect.Size.W) + 1
	}
}

func DrawSprite(sprite *model.Sprite, sectionName string, name string, X int32, Y int32) {
	section := sprite.GetSection(sectionName)
	rect := section.GetSprite(name)
	GlobalDisplay.DrawSpriteRect(sprite.Bitmap, rect, X, Y)
}

func DrawSpriteGradient(sprite *model.Sprite, sectionName string, name string, X int32, Y int32, sourceGradient model.Gradient, targetGradient model.Gradient, frameIndex int) {
	normalizeFrameIndex := int(frameIndex / 5)
	section := sprite.GetSection(sectionName)
	rect := section.GetSprite(name)
	GlobalDisplay.DrawSpriteRectGradient(sprite.Bitmap, rect, X, Y, sourceGradient, targetGradient, normalizeFrameIndex)
}

func DrawAnimation(sprite *model.Sprite, animationName string, frameIndex int, X int32, Y int32) {
	normalizeFrameIndex := int(frameIndex / 5)
	animation := sprite.GetAnimation(animationName)
	rects := sprite.GetAnimationRects(animation.Section)

	frames := animation.Frames

	rect := rects[frames[normalizeFrameIndex%len(frames)]]

	GlobalDisplay.DrawSpriteRect(sprite.Bitmap, rect, X, Y)
}

func (d *Display) FillRect(rect model.Rect, color []byte) {
	for y := 0; y < rect.Size.H; y++ {
		for x := 0; x < rect.Size.W; x++ {
			d.DrawPixel(int32(rect.Point.X+x), int32(rect.Point.Y+y), color)
		}
	}
}

func (d *Display) DrawSpriteRect(sprite *model.Bitmap, src model.Rect, destX, destY int32) {
	for sy := 0; sy < src.Size.H; sy++ {
		for sx := 0; sx < src.Size.W; sx++ {
			// Calculamos la posición real dentro del PNG original
			origX := src.Point.X + sx
			origY := src.Point.Y + sy

			// Seguridad: no leer fuera de la imagen original
			if origX < 0 || origX >= sprite.Size.W || origY < 0 || origY >= sprite.Size.H {
				continue
			}

			srcOff := (origY*sprite.Size.W + origX) * 4
			color := sprite.Pixels[srcOff : srcOff+4]

			// Transparencia
			if color[3] < 128 {
				continue
			}

			// Dibujar en la pantalla
			d.DrawPixel(destX+int32(sx), destY+int32(sy), color)
		}
	}
}

func (d *Display) DrawSpriteRectGradient(sprite *model.Bitmap, src model.Rect, destX, destY int32, sourceGradient model.Gradient, targetGradient model.Gradient, animationIndex int) {
	for sy := 0; sy < src.Size.H; sy++ {
		for sx := 0; sx < src.Size.W; sx++ {
			origX := src.Point.X + sx
			origY := src.Point.Y + sy

			if origX < 0 || origX >= sprite.Size.W || origY < 0 || origY >= sprite.Size.H {
				continue
			}

			srcOff := (origY*sprite.Size.W + sx + src.Point.X) * 4 // Asegúrate de sumar el offset X correctamente
			color := sprite.Pixels[srcOff : srcOff+4]

			if color[3] < 128 {
				continue
			}

			colorADibujar := ReplaceGradientColor(color, sourceGradient, targetGradient, animationIndex) // Por defecto el original

			d.DrawPixel(destX+int32(sx), destY+int32(sy), colorADibujar)
		}
	}
}

func ReplaceGradientColor(color []byte, sourceGradient model.Gradient, targetGradient model.Gradient, animationIndex int) []byte {
	gradientIndex := sourceGradient.GradientIndex(color)

	if gradientIndex >= 0 {
		return targetGradient[(gradientIndex+animationIndex)%len(targetGradient)].Byte()
	}
	return color
}
