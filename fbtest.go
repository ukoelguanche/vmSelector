package main

import (
	_ "image/png"
	"log"
	"time"

	"apodeiktikos.com/fbtest/drivers"
	"apodeiktikos.com/fbtest/loaders"
	"apodeiktikos.com/fbtest/model"
	"apodeiktikos.com/fbtest/util"
)

const targetFPS = 10
const frameDelay = time.Second / targetFPS

var gpuString string
var centinelVM *model.VM
var vms []model.VM

func DrawString(sprite *model.Sprite, text string, x, y int32, typography string) {
	cursorX := x

	letters := sprite.GetSection(typography)

	for _, char := range text {
		sChar := string(char)
		rect, ok := letters[sChar]
		if !ok {
			cursorX += 8
			continue
		}

		drivers.GlobalDisplay.DrawSpriteRect(sprite.Bitmap, rect, cursorX, y)
		cursorX += int32(rect.Size.W) + 1
	}
}

func DrawSprite(sprite *model.Sprite, sectionName string, name string, X int32, Y int32) {
	section := sprite.GetSection(sectionName)
	rect := section.GetSprite(name)
	drivers.GlobalDisplay.DrawSpriteRect(sprite.Bitmap, rect, X, Y)
}

func DrawAnimation(sprite *model.Sprite, animationName string, frameIndex int, X int32, Y int32) {
	animation := sprite.GetAnimation(animationName)
	rects := sprite.GetAnimationRects(animation.Section)

	frames := animation.Frames

	rect := rects[frames[frameIndex%len(frames)]]

	drivers.GlobalDisplay.DrawSpriteRect(sprite.Bitmap, rect, X, Y)
}

var hub *model.Sprite
var greenHill *model.Sprite
var sonic *model.Sprite

func GetVMsWithGPU(gpuString string, CentinelVM *model.VM) []model.VM {
	vms := model.GetVMs()
	if vms == nil {
		log.Fatal("Could not find any VMs")
	}

	var filtered []model.VM

	for _, vm := range vms.Data {
		if vm.HasSpecificGPU(gpuString) && vm.Name != CentinelVM.Name {
			filtered = append(filtered, vm)
		}
	}

	return filtered
}

func Init() {
	util.LoadContext()
	drivers.GlobalDisplay = drivers.InitDisplay(drivers.SW, drivers.SH, drivers.VW, drivers.VH)

	gpuString = util.ContextStorage.GpuString
	centinelVM = model.GetVMByName(util.ContextStorage.CentineVMName)
	vms = GetVMsWithGPU(gpuString, centinelVM)

	hub = loaders.LoadSprite("./resources/sprites/HUD.json")
	greenHill = loaders.LoadSprite("./resources/sprites/GreenHill.json")
	sonic = loaders.LoadSprite("./resources/sprites/Sonic.json")

}

func Loop(animationIndex int, selectedVMIndex int) {

	drivers.GlobalDisplay.Clear()

	colorFondo := []byte{255, 255, 255, 255}
	fondoRect := model.Rect{
		Point: model.Point{X: 0, Y: 0},
		Size:  model.Size{W: 320, H: 200},
	}

	drivers.GlobalDisplay.FillRect(fondoRect, colorFondo)
	DrawSprite(greenHill, "GreenHill", "background", 0, 0)

	DrawAnimation(greenHill, "flower1", animationIndex, 154, 90)
	DrawAnimation(greenHill, "flower2", animationIndex+15, -5, 115)
	DrawAnimation(greenHill, "flower2", animationIndex+7, 220, 115)
	DrawAnimation(greenHill, "flower2", animationIndex, 250, 115)

	const HUDX int32 = 130
	const HUDY int32 = 40

	for i, vm := range vms {
		var selectedOffset int32 = 0
		entryY := HUDY + int32(i*20)
		if i == selectedVMIndex {
			selectedOffset = 16
			DrawSprite(hub, "items", "emerald", HUDX, entryY+3)
		}
		DrawString(hub, vm.Name, HUDX+selectedOffset, entryY+1, "genesisLetters")
	}

	DrawAnimation(sonic, "sonic", animationIndex, 35, 128)

	drivers.GlobalDisplay.Present()
}

func main() {
	Init()

	defer drivers.GlobalDisplay.Close()
	var animationIndex = 0
	var selectedVMIndex = 0

	for {
		dx, dy, quit := drivers.GlobalDisplay.GetInput()
		if quit {
			continue
		}

		if dx > 0 || dy > 0 {
			selectedVMIndex += 1
		}
		if dx < 0 || dy < 0 {
			selectedVMIndex -= 1
		}

		start := time.Now()
		Loop(animationIndex, selectedVMIndex)

		animationIndex++

		elapsed := time.Since(start)
		// log.Printf("elapsed %s", elapsed)
		if elapsed < frameDelay {
			time.Sleep(frameDelay - elapsed)
		}
	}
}
