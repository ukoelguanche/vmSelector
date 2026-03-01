package main

import (
	"os"
)

const (
	vW, vH = 320, 200  // Resolución lógica (Retro)
	sW, sH = 1024, 768 // Resolución real del FB (Ajusta según tu GRUB)
)

func main() {
	// 1. Abrimos el dispositivo del Framebuffer
	fbFile, err := os.OpenFile("/dev/fb0", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer fbFile.Close()

	// 2. Creamos un buffer en RAM para toda la pantalla (1024 * 768 * 4 bytes)
	// Dibujar primero en RAM y luego volcar es mucho más rápido
	screen := make([]byte, sW*sH*4)

	// 3. Vamos a dibujar algo: un degradado de fondo usando tus píxeles gordos
	for y := 0; y < vH; y++ {
		for x := 0; x < vW; x++ {
			// Color dinámico basado en coordenadas (formato BGRA)
			// Blue, Green, Red, Alpha
			color := []byte{byte(x % 255), byte(y % 255), 100, 255}
			drawFatPixel(screen, x, y, color)
		}
	}

	// 4. Volcamos todo el buffer al hardware de una sola vez
	fbFile.Write(screen)
}

// Tu función (corregida para escribir en el slice 'screen' en lugar de directo al archivo)
func drawFatPixel(fb []byte, x, y int, color []byte) {
	scaleX := sW / vW
	scaleY := sH / vH

	for py := 0; py < scaleY; py++ {
		for px := 0; px < scaleX; px++ {
			realX := x*scaleX + px
			realY := y*scaleY + py
			
			// Seguridad: no escribir fuera del array
			if realX < sW && realY < sH {
				offset := (realY*sW + realX) * 4
				copy(fb[offset:offset+4], color)
			}
		}
	}
}
