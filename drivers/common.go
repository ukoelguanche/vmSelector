package drivers

import "apodeiktikos.com/fbtest/model"

const (
	VW, VH = 320, 200
	SW, SH = 640, 480
)

var GlobalDisplay *Display

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
			if origX < 0 || origX >= sprite.W || origY < 0 || origY >= sprite.H {
				continue
			}

			srcOff := (origY*sprite.W + origX) * 4
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
