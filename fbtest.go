package main

import (
	"fmt"
	_ "image/png"
	"time"

	"apodeiktikos.com/fbtest/drivers"
	"apodeiktikos.com/fbtest/loaders"
	"apodeiktikos.com/fbtest/model"
)

const targetFPS = 30
const frameDelay = time.Second / targetFPS

var fuenteMapa map[string]map[string]model.Rect

func DrawString(display *drivers.Display, sprite *model.Sprite, text string, x, y int32) {
	cursorX := x
	letras, ok := fuenteMapa["boldLetters"]
	if !ok {
		return // No existe la sección "letters" en el JSON
	}

	for _, char := range text {
		sChar := string(char)
		rect, ok := letras[sChar]
		if !ok {
			cursorX += 8
			continue
		}

		display.DrawSpriteRect(sprite, rect, cursorX, y)
		cursorX += int32(rect.Size.W) + 1
	}
}

func DrawSprite(display *drivers.Display, sprite *model.Sprite, sectionName string, name string) {

	section, ok := fuenteMapa[sectionName]
	if !ok {
		return // No existe la sección "letters" en el JSON
	}

	rect, ok := section[name]
	if !ok {
		return
	}

	display.DrawSpriteRect(sprite, rect, 100, 100)
}

func main() {
	display := drivers.InitDisplay(drivers.SW, drivers.SH, drivers.VW, drivers.VH)
	defer display.Close()

	miSprite, err := loaders.LoadPNG("./resources/sprites/HUD.png")
	if err != nil {
		fmt.Println("Error cargando sprite:", err)
		return
	}

	err2 := loaders.LoadJSON("./resources/sprites/HUD.json", &fuenteMapa)
	if err2 != nil {
		fmt.Println("Error cargando json:", err)
		return
	}

	if miSprite == nil {
		fmt.Println("Error cargando sprite")
	}

	var x, y int32 = 50, 50

	for {
		start := time.Now()

		dx, dy, quit := display.GetInput()
		if quit {
			break
		}
		x += dx
		y += dy

		display.Clear()

		//colorFondo := []byte{40, 40, 40, 255} // Gris
		colorFondo := []byte{0, 0, 180, 255} // Azul
		fondoRect := model.Rect{
			Point: model.Point{X: 0, Y: 0},
			Size:  model.Size{W: 320, H: 200},
		}

		display.FillRect(fondoRect, colorFondo)

		DrawSprite(display, miSprite, "panel", "top")

		DrawString(display, miSprite, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", 10, 10)
		DrawString(display, miSprite, "abcdefghijklmnopqrstuvwxyz", 10, 20)
		DrawString(display, miSprite, "0123456789", 10, 30)
		DrawString(display, miSprite, "apodeiktikos", 133, 105)

		display.Present()
		//time.Sleep(16 * time.Millisecond)

		elapsed := time.Since(start)
		// log.Println("Elapsed time:", elapsed)
		if elapsed < frameDelay {
			time.Sleep(frameDelay - elapsed) // Dormimos el resto hasta llegar a los 33.3ms
		}
	}
}
