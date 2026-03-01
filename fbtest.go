package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"time"

	"apodeiktikos.com/fbtest/model"
)

func LoadPNG(path string) (*model.Sprite, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	pixels := make([]byte, w*h*4)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			offset := (y*w + x) * 4
			// Convertimos de 16-bit a 8-bit que devuelve Go
			pixels[offset] = byte(r >> 8)
			pixels[offset+1] = byte(g >> 8)
			pixels[offset+2] = byte(b >> 8)
			pixels[offset+3] = byte(a >> 8)
		}
	}
	return &model.Sprite{W: w, H: h, Pixels: pixels}, nil
}

const (
	vW, vH = 320, 200
	sW, sH = 640, 480
)

const targetFPS = 30
const frameDelay = time.Second / targetFPS

var fuenteMapa map[string]map[string]model.Rect

func LoadFontConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &fuenteMapa)
}

func DrawString(display *Display, sprite *model.Sprite, text string, x, y int32) {
	cursorX := x
	letras, ok := fuenteMapa["letters"]
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

func DrawSprite(display *Display, sprite *model.Sprite, sectionName string, name string) {

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

func (d *Display) FillRect(rect model.Rect, color []byte) {
	for y := 0; y < rect.Size.H; y++ {
		for x := 0; x < rect.Size.W; x++ {
			// Usamos tu lógica de DrawPixel para que se vea
			// bien tanto en Mac como con el "píxel gordo" de Linux
			d.DrawPixel(int32(rect.Point.X+x), int32(rect.Point.Y+y), color)
		}
	}
}

func main() {
	display := InitDisplay(sW, sH, vW, vH)
	defer display.Close()

	miSprite, err := LoadPNG("./resources/sprites/HUD.png")
	if err != nil {
		fmt.Println("Error cargando sprite:", err)
		return
	}

	err2 := LoadFontConfig("./resources/sprites/HUD.json")
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

		DrawString(display, miSprite, "abcdefghijklmnopqrstuvwxyz", 0, 0)
		//DrawString(display, miSprite, "0123456789  ", 0, 9)
		//DrawString(display, miSprite, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", 0, 9)

		DrawSprite(display, miSprite, "panel", "bottom")

		/*
			str := "H"
			chr := str[0]
			idx := int(chr) - 65

			// 3. Calculamos la posición (suponiendo que cada letra mide, por ejemplo, 16px)
			const charWidth = 9
			xSource := idx * charWidth

			source := Rect{X: xSource, Y: 0, W: 8, H: 8}
			display.DrawSpriteRect(miSprite, source, x, y)
		*/

		display.Present()
		//time.Sleep(16 * time.Millisecond)

		elapsed := time.Since(start) // ¿Cuánto tiempo hemos gastado trabajando?
		// log.Println("Elapsed time:", elapsed)
		if elapsed < frameDelay {
			time.Sleep(frameDelay - elapsed) // Dormimos el resto hasta llegar a los 33.3ms
		}
	}
}
