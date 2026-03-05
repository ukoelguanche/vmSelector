package drivers

import "apodeiktikos.com/fbtest/model"

const (
	VW, VH = 320, 200
)

var GlobalDisplay *Display

func DrawSpriteFrame(sprite *model.SpriteInstance) {
	GlobalDisplay.DrawSpriteRect(sprite, sprite.CurrentFrame(), sprite.Position)
}

func DrawText(text *model.Text) {
	cursorX := text.Position.X

	letters := text.Sprite.Frames
	characters := text.Sprite.Characters

	for _, char := range text.Text {
		sChar := string(char)
		rect := letters[characters[sChar]]

		GlobalDisplay.DrawSpriteRect(text, rect, model.Point{X: cursorX, Y: text.Position.Y})
		cursorX += rect.Size.W + 1
	}
}

func (d *Display) DrawSpriteRect(sprite model.Renderable, rect model.Rect, position model.Point) {
	bitmap := sprite.GetBitmap()
	for sy := 0; sy < int(rect.Size.H); sy++ {
		for sx := 0; sx < int(rect.Size.W); sx++ {
			// Calculate original position inside bitmap
			origX := rect.Point.X + float64(sx)
			origY := rect.Point.Y + float64(sy)

			// Avoid drawing outside bounds
			if origX < 0 || origX >= float64(bitmap.W) || origY < 0 || origY >= float64(bitmap.H) {
				continue
			}

			srcOff := int((origY*float64(bitmap.W) + origX) * 4)
			color := bitmap.Pixels[srcOff : srcOff+4]

			// Skip transparencies
			if color[3] < 128 {
				continue
			}

			finalColor := sprite.ProcessColor(color)

			d.DrawPixel(int32(position.X)+int32(sx), int32(position.Y)+int32(sy), finalColor)
		}
	}
}
