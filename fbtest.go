package main

import (
	_ "image/png"
	"log"
	"time"

	"apodeiktikos.com/fbtest/drivers"
	"apodeiktikos.com/fbtest/loaders"
	"apodeiktikos.com/fbtest/model"
)

const targetFPS = 10
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

		display.DrawSpriteRect(sprite.Bitmap, rect, cursorX, y)
		cursorX += int32(rect.Size.W) + 1
	}
}

func DrawSprite(sprite *model.Sprite, display *drivers.Display, sectionName string, name string, X int32, Y int32) {
	section := sprite.GetSection(sectionName)
	rect := section.GetSprite(name)
	display.DrawSpriteRect(sprite.Bitmap, rect, X, Y)
}

func DrawAnimation(sprite *model.Sprite, display *drivers.Display, animationName string, frameIndex int, X int32, Y int32) {
	animation := sprite.GetAnimation(animationName)
	rects := sprite.GetAnimationRects(animation.Section)

	frames := animation.Frames

	rect := rects[frames[frameIndex%len(frames)]]

	display.DrawSpriteRect(sprite.Bitmap, rect, X, Y)
}

var hudSprite *model.Sprite
var rossi *model.Sprite
var sonic *model.Sprite

func Init() {
	hudSprite = loaders.LoadSprite("./resources/sprites/HUD.json")
	rossi = loaders.LoadSprite("./resources/sprites/rossi.json")
	sonic = loaders.LoadSprite("./resources/sprites/Sonic.json")
	return
}

func main() {
	Init()

	// ToDo: Convert display to global variable
	display := drivers.InitDisplay(drivers.SW, drivers.SH, drivers.VW, drivers.VH)
	defer display.Close()

	var x, y int32 = 50, 50

	var animationIndex = 0
	for {
		start := time.Now()

		dx, dy, quit := display.GetInput()
		if quit {
			break
		}

		display.Clear()

		colorFondo := []byte{255, 255, 255, 255}
		fondoRect := model.Rect{
			Point: model.Point{X: 0, Y: 0},
			Size:  model.Size{W: 320, H: 200},
		}
		log.Printf("%d %d", dx, dy)
		display.FillRect(fondoRect, colorFondo)

		DrawSprite(hudSprite, display, "panel", "top", 80, 20)
		DrawSprite(hudSprite, display, "panel", "center", 80, 38)
		DrawSprite(hudSprite, display, "panel", "center", 80, 56)
		DrawSprite(hudSprite, display, "panel", "bottom", 80, 72)

		DrawAnimation(sonic, display, "sonic", animationIndex, 35, 132)

		DrawString(display, hudSprite, "Apodeiktikos", 86, 26)

		display.Present()

		elapsed := time.Since(start)
		if elapsed < frameDelay {
			time.Sleep(frameDelay - elapsed)
		}
		animationIndex = animationIndex + 1
	}
}
