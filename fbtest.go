package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	vW, vH = 320, 200  // Resolución lógica (Retro)
	sW, sH = 1024, 768 // Resolución real (Ventana/FB)
)

func main() {
	// 1. Inicializar SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Error init: %s\n", err)
		return
	}
	defer sdl.Quit()

	// 2. Crear Ventana
	window, err := sdl.CreateWindow("Mi Selector Retro",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(sW), int32(sH), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ventana: %s\n", err)
		return
	}
	defer window.Destroy()

	// 3. Crear el Renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return
	}
	defer renderer.Destroy()

	// 4. Crear la Textura "Retro"
	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING, int32(vW), int32(vH))

	var x, y int32 = 160, 100

	running := true
	for running {
		// --- GESTIONAR EVENTOS ---
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					switch t.Keysym.Sym {
					case sdl.K_UP:
						y -= 4
					case sdl.K_DOWN:
						y += 4
					case sdl.K_LEFT:
						x -= 4
					case sdl.K_RIGHT:
						x += 4
					case sdl.K_ESCAPE:
						running = false
					}
				}
			} // <-- AQUÍ faltaba cerrar el switch de eventos
		} // <-- AQUÍ faltaba cerrar el for de PollEvent

		// --- DIBUJAR ---
		// Limpiamos el fondo del renderer (opcional si la textura ocupa todo)
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// Preparamos el array de píxeles (vW * vH * 4 bytes por píxel)
		pixels := make([]byte, vW*vH*4)

		// Dibujamos nuestro cuadrado blanco de 10x10 en la posición x, y
		for cy := int32(0); cy < 10; cy++ {
			for cx := int32(0); cx < 10; cx++ {
				targetX, targetY := x+cx, y+cy
				// Evitamos salirnos de los límites de la textura de 320x200
				if targetX >= 0 && targetX < vW && targetY >= 0 && targetY < vH {
					offset := (targetY*vW + targetX) * 4
					pixels[offset] = 255   // R
					pixels[offset+1] = 255 // G
					pixels[offset+2] = 255 // B
					pixels[offset+3] = 255 // A
				}
			}
		}

		// Actualizamos la textura con los nuevos píxeles
		// texture.Update(nil, pixels, vW*4)
		texture.Update(nil, unsafe.Pointer(&pixels[0]), vW*4)

		// Copiamos la textura pequeña (320x200) al renderer grande (1024x768)
		// SDL hará el escalado automáticamente por nosotros
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		sdl.Delay(16) // Equivale a unos 60 FPS
	}
}
