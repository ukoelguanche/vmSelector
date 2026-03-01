package main

import (
	_ "image/png"
	"time"

	"apodeiktikos.com/fbtest/drivers"
	"apodeiktikos.com/fbtest/loaders"
	"apodeiktikos.com/fbtest/model"
)

const targetFPS = 30
const frameDelay = time.Second / targetFPS

func DrawString(display *drivers.Display, sprite *model.Sprite, text string, x, y int32) {
	cursorX := x

	letters := sprite.GetSection("boldLetters")

	for _, char := range text {
		sChar := string(char)
		rect, ok := letters[sChar]
		if !ok {
			cursorX += 8
			continue
		}

		display.DrawSpriteRect2(sprite.Bitmap, rect, cursorX, y)
		cursorX += int32(rect.Size.W) + 1
	}
}

func DrawSprite(sprite *model.Sprite, display *drivers.Display, sectionName string, name string, X int32, Y int32) {
	section := sprite.GetSection(sectionName)
	rect := section.GetSprite(name)
	display.DrawSpriteRect2(sprite.Bitmap, rect, X, Y)
}

var hudSprite *model.Sprite

func Init() {
	hudSprite = loaders.LoadSprite("./resources/sprites/HUD.json")
	return
}

func main() {
	Init()

	display := drivers.InitDisplay(drivers.SW, drivers.SH, drivers.VW, drivers.VH)
	defer display.Close()

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

		DrawSprite(hudSprite, display, "panel", "top", 80, 20)
		DrawSprite(hudSprite, display, "panel", "center", 80, 38)
		DrawSprite(hudSprite, display, "panel", "center", 80, 56)
		DrawSprite(hudSprite, display, "panel", "bottom", 80, 72)

		DrawString(display, hudSprite, "A P O D E I K T I K O S", 86, 26)

		display.Present()

		elapsed := time.Since(start)
		if elapsed < frameDelay {
			time.Sleep(frameDelay - elapsed) // Dormimos el resto hasta llegar a los 33.3ms
		}
	}
}
