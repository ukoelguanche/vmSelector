package drivers

import "apodeiktikos.com/fbtest/model"

const (
	VW, VH = 320, 200
	SW, SH = 1280, 720
)

var GlobalDisplay *Display

func DrawText(text *model.Text) {
	cursorX := text.Position.X

	letters := text.Sprite.Frames
	characters := text.Sprite.Characters

	for _, char := range text.Text {
		sChar := string(char)
		rect := letters[characters[sChar]]

		GlobalDisplay.DrawSpriteRect(text.Sprite, rect, cursorX, text.Position.Y)
		cursorX += int32(rect.Size.W) + 1
	}

}

/*
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
*/

func DrawAnimation(sprite *model.SpriteInstance) {
	GlobalDisplay.DrawSpriteRect(sprite.Sprite, sprite.CurrentFrame(), sprite.Position.X, sprite.Position.Y)
}

func (d *Display) FillRect(rect model.Rect, color []byte) {
	for y := 0; y < int(rect.Size.H); y++ {
		for x := 0; x < int(rect.Size.W); x++ {
			d.DrawPixel(rect.Point.X+int32(x), rect.Point.Y+int32(y), color)
		}
	}
}

func (d *Display) DrawSpriteRect(sprite *model.Sprite, src model.Rect, destX, destY int32) {
	bitmap := sprite.Bitmap
	for sy := 0; sy < int(src.Size.H); sy++ {
		for sx := 0; sx < int(src.Size.W); sx++ {
			// Calculamos la posición real dentro del PNG original
			origX := src.Point.X + int32(sx)
			origY := src.Point.Y + int32(sy)

			// Seguridad: no leer fuera de la imagen original
			if origX < 0 || origX >= bitmap.Size.W || origY < 0 || origY >= bitmap.Size.H {
				continue
			}

			srcOff := (origY*bitmap.Size.W + origX) * 4
			color := bitmap.Pixels[srcOff : srcOff+4]

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
	for sy := 0; sy < int(src.Size.H); sy++ {
		for sx := 0; sx < int(src.Size.W); sx++ {
			origX := src.Point.X + int32(sx)
			origY := src.Point.Y + int32(sy)

			if origX < 0 || origX >= sprite.Size.W || origY < 0 || origY >= sprite.Size.H {
				continue
			}

			srcOff := (origY*sprite.Size.W + int32(sx) + src.Point.X) * 4 // Asegúrate de sumar el offset X correctamente
			color := sprite.Pixels[srcOff : srcOff+4]

			if color[3] < 128 {
				continue
			}

			//colorADibujar := ReplaceGradientColor(color, sourceGradient, targetGradient, animationIndex) // Por defecto el original
			colorADibujar := color

			d.DrawPixel(destX+int32(sx), destY+int32(sy), colorADibujar)
		}
	}
}

/*
func ReplaceGradientColor(color []byte, sourceGradient model.Gradient, targetGradient model.Gradient, animationIndex int) []byte {
	gradientIndex := sourceGradient.GradientIndex(color)

	if gradientIndex >= 0 {
		return targetGradient[(gradientIndex+animationIndex)%len(targetGradient)].Byte()
	}
	return color
}

*/
