package main

import "os"

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
	// Aquí va tu lógica de "Píxel Gordo" manual
	scaleX, scaleY := sW/vW, sH/vH
	for py := 0; py < scaleY; py++ {
		for px := 0; px < scaleX; px++ {
			rx, ry := int(vx)*scaleX+px, int(vy)*scaleY+py
			if rx < sW && ry < sH {
				offset := (ry*sW + rx) * 4
				copy(d.pixels[offset:offset+4], c)
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

func (d *Display) DrawSprite(sprite *Sprite, x, y int32) {
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
