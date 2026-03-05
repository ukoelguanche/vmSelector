package drivers

import "apodeiktikos.com/fbtest/model"

const (
	VW, VH = 320, 200
)

var GlobalDisplay *Display

func DrawSpriteFrame(sprite *model.SpriteInstance) {
	GlobalDisplay.DrawSpriteRect(sprite, sprite.CurrentFrame(), sprite.Position.X, sprite.Position.Y)
}

func DrawText(text *model.Text) {
	cursorX := text.Position.X

	letters := text.Sprite.Frames
	characters := text.Sprite.Characters

	for _, char := range text.Text {
		sChar := string(char)
		rect := letters[characters[sChar]]

		GlobalDisplay.DrawSpriteRect(text, rect, cursorX, text.Position.Y)
		cursorX += int32(rect.Size.W) + 1
	}
}

func (d *Display) DrawSpriteRect(sprite model.Renderable, rect model.Rect, destX, destY int32) {
	bitmap := sprite.GetBitmap()
	for sy := 0; sy < int(rect.Size.H); sy++ {
		for sx := 0; sx < int(rect.Size.W); sx++ {
			// Calculate original position inside bitmap
			origX := rect.Point.X + int32(sx)
			origY := rect.Point.Y + int32(sy)

			// Avoid drawing outside bounds
			if origX < 0 || origX >= bitmap.Size.W || origY < 0 || origY >= bitmap.Size.H {
				continue
			}

			srcOff := (origY*bitmap.Size.W + origX) * 4
			color := bitmap.Pixels[srcOff : srcOff+4]

			// Skip transparencies
			if color[3] < 128 {
				continue
			}

			finalColor := sprite.ProcessColor(color)

			d.DrawPixel(destX+int32(sx), destY+int32(sy), finalColor)
		}
	}
}
