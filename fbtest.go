package main

import (
	"time"
)

const (
	vW, vH = 320, 200
	sW, sH = 640, 480
)

const targetFPS = 30
const frameDelay = time.Second / targetFPS

func main() {
	display := InitDisplay(sW, sH, vW, vH)
	defer display.Close()

	var x, y int32 = 160, 100

	var t int32 = 0

	for {
		start := time.Now()

		dx, dy, quit := display.GetInput()
		if quit {
			break
		}
		x += dx
		y += dy

		display.Clear()

		color := []byte{255, 255, 255, 255} // Blanco
		for cy := int32(0); cy < 10; cy++ {
			for cx := int32(0); cx < 10; cx++ {
				display.DrawPixel(x+cx+t, y+cy, color)
			}
		}
		t = (t + 2) % 100

		display.Present()
		//time.Sleep(16 * time.Millisecond)

		elapsed := time.Since(start) // ¿Cuánto tiempo hemos gastado trabajando?
		if elapsed < frameDelay {
			time.Sleep(frameDelay - elapsed) // Dormimos el resto hasta llegar a los 33.3ms
		}
	}
}
