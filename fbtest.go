package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"
	"time"

	"apodeiktikos.com/fbtest/drivers"
	"apodeiktikos.com/fbtest/loaders"
	"apodeiktikos.com/fbtest/model"
	"apodeiktikos.com/fbtest/util"
)

const targetFPS = 60
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

func DrawSpriteGradient(sprite *model.Sprite, sectionName string, name string, X int32, Y int32, sourceGradient model.Gradient, targetGradient model.Gradient, frameIndex int) {
	normalizeFrameIndex := int(frameIndex / 5)
	section := sprite.GetSection(sectionName)
	rect := section.GetSprite(name)
	drivers.GlobalDisplay.DrawSpriteRectGradient(sprite.Bitmap, rect, X, Y, sourceGradient, targetGradient, normalizeFrameIndex)
}

func DrawAnimation(sprite *model.Sprite, animationName string, frameIndex int, X int32, Y int32) {
	normalizeFrameIndex := int(frameIndex / 5)
	animation := sprite.GetAnimation(animationName)
	rects := sprite.GetAnimationRects(animation.Section)

	frames := animation.Frames

	rect := rects[frames[normalizeFrameIndex%len(frames)]]

	drivers.GlobalDisplay.DrawSpriteRect(sprite.Bitmap, rect, X, Y)
}

var hub *model.Sprite
var greenHill *model.Sprite
var greenHillBack *model.Sprite
var sonic *model.Sprite

func GetVMsWithGPU(gpuString string, centinelVM *model.VM) []model.VM {
	vms := model.GetVMs()
	if vms == nil {
		log.Fatal("Could not find any VMs")
	}

	var filtered []model.VM

	for _, vm := range vms.Data {
		if vm.HasSpecificGPU(gpuString) && vm.Name != centinelVM.Name {
			filtered = append(filtered, vm)
		}
	}

	filtered = append(filtered, *centinelVM)

	return filtered
}

func SwitchToVM(centinelVM *model.VM, targetVM model.VM) {
	if centinelVM.VMID == targetVM.VMID {
		model.SetVMDescription(centinelVM, "power_off")
	} else {
		model.SetVMDescription(centinelVM, fmt.Sprintf("target_vm_id %d", targetVM.VMID))
	}
	model.PowerOffVM(centinelVM)
}

func Init() {
	util.LoadContext()
	drivers.GlobalDisplay = drivers.InitDisplay(drivers.SW, drivers.SH, drivers.VW, drivers.VH)

	gpuString = util.ContextStorage.GpuString
	centinelVM = model.GetVMByName(util.ContextStorage.CentineVMName)
	vms = GetVMsWithGPU(gpuString, centinelVM)

	hub = loaders.LoadSprite("./resources/sprites/HUD.json")
	greenHill = loaders.LoadSprite("./resources/sprites/GreenHill.json")
	greenHillBack = loaders.LoadSprite("./resources/sprites/GreenHillBack.json")
	sonic = loaders.LoadSprite("./resources/sprites/Sonic.json")

}

func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}
func RenderHUD(animationIndex int, selectedVMIndex int) {
	const initialPos = 220
	animationPercent := min(float64(animationIndex), initialPos) / initialPos
	ease := EaseInOutCubic(animationPercent)

	xOffset := initialPos - int32(initialPos*ease)

	for i := 0; i < 13; i++ {
		DrawSprite(hub, "items", "zigZagBG", 90+xOffset, int32(i*16))
	}

	var HUDX int32 = 130 + xOffset
	var HUDY int32 = 40

	var selectedOffset int32 = 0
	for i, vm := range vms {
		entryY := HUDY + int32(i*20)
		text := vm.Name

		if i == len(vms)-1 {
			entryY = entryY + 20
			text = "Power Down"
		}

		if i == selectedVMIndex {
			selectedOffset = 4
			DrawAnimation(hub, "ring", animationIndex, HUDX-15, entryY+1)
		} else {
			selectedOffset = 0
		}

		DrawString(hub, text, HUDX+selectedOffset, entryY+1, "genesisLetters")
	}
}

func Loop(animationIndex int, selectedVMIndex int, endLoop bool) {

	drivers.GlobalDisplay.Clear()

	colorFondo := []byte{255, 255, 255, 255}
	fondoRect := model.Rect{
		Point: model.Point{X: 0, Y: 0},
		Size:  model.Size{W: 320, H: 200},
	}

	drivers.GlobalDisplay.FillRect(fondoRect, colorFondo)

	sourceGradient := []model.Color{
		{R: 221, G: 119, B: 221, A: 255},
		{R: 187, G: 85, B: 187, A: 255},
		{R: 153, G: 51, B: 153, A: 255},
		{R: 119, G: 17, B: 119, A: 255},
	}

	targetGradient := []model.Color{
		{R: 151, G: 179, B: 246, A: 255},
		{R: 115, G: 143, B: 245, A: 255},
		{R: 115, G: 143, B: 177, A: 255},
		{R: 187, G: 215, B: 249, A: 255},
	}

	DrawSprite(greenHillBack, "GreenHillBack", "layer6", 0-int32(float64(animationIndex)*0.2), 0)
	DrawSprite(greenHillBack, "GreenHillBack", "layer5", 0-int32(float64(animationIndex)*0.1), 32)
	DrawSprite(greenHillBack, "GreenHillBack", "layer4", 0-int32(float64(animationIndex)*0.05), 48)
	DrawSprite(greenHillBack, "GreenHillBack", "layer3", 0, 64)
	DrawSpriteGradient(greenHillBack, "GreenHillBack", "layer2", 0, 112, sourceGradient, targetGradient, animationIndex)
	DrawSpriteGradient(greenHillBack, "GreenHillBack", "layer1", 0, 152, sourceGradient, targetGradient, animationIndex)

	DrawSprite(greenHill, "GreenHill", "background", 0, 0)

	DrawAnimation(greenHill, "flower1", animationIndex, 154, 90)
	DrawAnimation(greenHill, "flower2", animationIndex+15, -5, 115)
	DrawAnimation(greenHill, "flower2", animationIndex+7, 220, 115)
	DrawAnimation(greenHill, "flower2", animationIndex, 250, 115)

	RenderHUD(animationIndex, selectedVMIndex)

	DrawAnimation(sonic, "sonic", animationIndex, 35, 128)

	drivers.GlobalDisplay.Present()
}

func main() {
	Init()

	defer drivers.GlobalDisplay.Close()
	var animationIndex = 0
	var selectedVMIndex = 0
	var endLoop = false

	for {
		dx, dy, quit, enter := drivers.GlobalDisplay.GetInput()
		if quit {
			break
		}
		if enter {
			endLoop = true
			SwitchToVM(centinelVM, vms[selectedVMIndex])
		}

		if (dx > 0 || dy > 0) && selectedVMIndex < len(vms) {
			selectedVMIndex += 1
		}
		if (dx < 0 || dy < 0) && selectedVMIndex > 0 {
			selectedVMIndex -= 1
		}

		start := time.Now()

		Loop(animationIndex, selectedVMIndex, endLoop)

		animationIndex += 1

		elapsed := time.Since(start)

		if elapsed < frameDelay {
			time.Sleep(frameDelay - elapsed)
		}
	}
}
