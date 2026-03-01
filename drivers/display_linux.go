package drivers

import (
	"os"

	"apodeiktikos.com/fbtest/model"
)

type Display struct {
	file   *os.File
	pixels []byte
}

func InitDisplay(sw, sh, vw, vh int) *Display {
	f, _ := os.OpenFile("/dev/fb0", os.O_RDWR, 0)
	return &Display{
		file:   f,
		pixels: make([]byte, sw*sh*4),
	}
}

func (d *Display) DrawPixel(vx, vy int32, c []byte) {
	if vx < 0 || vx >= VW || vy < 0 || vy >= VH {
		return
	}

	scaleX, scaleY := SW/VW, SH/VH

	// Extraemos los componentes del color original (RGBA)
	r, g, b, a := c[0], c[1], c[2], c[3]

	for py := 0; py < scaleY; py++ {
		for px := 0; px < scaleX; px++ {
			rx, ry := int(vx)*scaleX+px, int(vy)*scaleY+py

			if rx >= 0 && rx < SW && ry >= 0 && ry < SH {
				offset := (ry*SW + rx) * 4

				// AQUÍ ESTÁ EL CAMBIO: Escribimos en orden B, G, R, A
				d.pixels[offset] = b   // Azul primero
				d.pixels[offset+1] = g // Verde igual
				d.pixels[offset+2] = r // Rojo al final
				d.pixels[offset+3] = a // Alpha
			}
		}
	}
}

func (d *Display) Clear() {
	for i := range d.pixels {
		d.pixels[i] = 0
	}
}

func (d *Display) Present() {
	d.file.WriteAt(d.pixels, 0)
}

func (d *Display) GetInput() (int32, int32, bool) {
	// Por ahora estático, luego leeremos /dev/input/event0
	return 0, 0, false
}

func (d *Display) Close() { d.file.Close() }

func (d *Display) DrawSprite(sprite *model.Sprite, x, y int32) {
	for sy := 0; sy < sprite.H; sy++ {
		for sx := 0; sx < sprite.W; sx++ {
			srcOff := (sy*sprite.W + sx) * 4
			color := sprite.Pixels[srcOff : srcOff+4]

			// Si el píxel es transparente (Alpha < 128), no lo dibujamos
			if color[3] < 128 {
				continue
			}

			// Dibujamos el píxel usando nuestra lógica de "píxel gordo"
			d.DrawPixel(x+int32(sx), y+int32(sy), color)
		}
	}
}
