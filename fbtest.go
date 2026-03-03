package main

import (
	_ "image/png"
	"time"

	"apodeiktikos.com/fbtest/drivers"
	"apodeiktikos.com/fbtest/loaders"
	"apodeiktikos.com/fbtest/model"
	//"apodeiktikos.com/fbtest/sprites"
	"apodeiktikos.com/fbtest/util"
)

const targetFPS = 30
const frameDelay = time.Second / targetFPS

var gpuString string
var centinelVM *model.VM
var vms []model.VM

var sprites model.Sprites
var spriteInstances []*model.SpriteInstance
var texts []*model.Text

func Init() {
	util.LoadContext()
	drivers.GlobalDisplay = drivers.InitDisplay(drivers.SW, drivers.SH, drivers.VW, drivers.VH)

	loaders.LoadSprites("./resources/sprites/Sprites.json", &sprites)

	spriteInstances = make([]*model.SpriteInstance, 0)
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer1", "idle", model.Point{X: 0, Y: 0}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer2", "idle", model.Point{X: 0, Y: 32}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer3", "idle", model.Point{X: 0, Y: 48}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer4", "idle", model.Point{X: 0, Y: 64}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer5", "idle", model.Point{X: 0, Y: 112}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer6", "idle", model.Point{X: 0, Y: 152}))

	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower1", "idle", model.Point{X: 154, Y: 90}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower2", "idle", model.Point{X: -5, Y: 115}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower2", "idle", model.Point{X: 220, Y: 115}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower2", "idle", model.Point{X: 250, Y: 115}))

	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Sonic", "idle", model.Point{X: 35, Y: 130}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillForeground", "idle", model.Point{X: 0, Y: 0}))
	for y := 0; y < 13; y++ {
		spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "ZigZag", "idle", model.Point{X: 100, Y: int32(y * 16)}))
	}
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Ring", "idle", model.Point{X: 130, Y: 50}))

	texts = make([]*model.Text, 0)

	gpuString = util.ContextStorage.GpuString
	centinelVM = model.GetVMByName(util.ContextStorage.CentineVMName)
	vms = model.GetVMsWithGPU(gpuString, centinelVM)
	for i, vm := range vms {
		texts = append(texts, &model.Text{Sprite: sprites.Sprites["GenesisLetters"], Text: vm.Name, Position: model.Point{X: 130, Y: int32(i)*20 + 30}})
	}

	return

}

/*
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
		drivers.DrawSprite(hud, "items", "zigZagBG", 90+xOffset, int32(i*16))
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
			drivers.DrawAnimation(hud, "ring", animationIndex, HUDX-15, entryY+1)
		} else {
			selectedOffset = 0
		}

		drivers.DrawString(hud, text, HUDX+selectedOffset, entryY+1, "genesisLetters")
	}
}
*/

func Loop(animationIndex int, selectedVMIndex int, endLoop bool) {
	drivers.GlobalDisplay.Clear()

	/*
		colorFondo := []byte{255, 255, 255, 255}
		fondoRect := model.Rect{
			Point: model.Point{X: 0, Y: 0},
			Size:  model.Size{W: 320, H: 200},
		}

		drivers.GlobalDisplay.FillRect(fondoRect, colorFondo)
	*/

	for _, spriteInstance := range spriteInstances {
		drivers.DrawAnimation(spriteInstance)
		spriteInstance.NextFrame()
	}

	for _, text := range texts {
		drivers.DrawText(text)
	}

	/*

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

		drivers.DrawSprite(greenHillBack, "GreenHillBack", "layer6", 0-int32(float64(animationIndex)*0.2), 0)
		drivers.DrawSprite(greenHillBack, "GreenHillBack", "layer5", 0-int32(float64(animationIndex)*0.1), 32)
		drivers.DrawSprite(greenHillBack, "GreenHillBack", "layer4", 0-int32(float64(animationIndex)*0.05), 48)
		drivers.DrawSprite(greenHillBack, "GreenHillBack", "layer3", 0, 64)
		drivers.DrawSpriteGradient(greenHillBack, "GreenHillBack", "layer2", 0, 112, sourceGradient, targetGradient, animationIndex)
		drivers.DrawSpriteGradient(greenHillBack, "GreenHillBack", "layer1", 0, 152, sourceGradient, targetGradient, animationIndex)

		drivers.DrawSprite(greenHill, "GreenHill", "background", 0, 0)

		drivers.DrawAnimation(greenHill, "flower1", animationIndex, 154, 90)
		drivers.DrawAnimation(greenHill, "flower2", animationIndex+15, -5, 115)
		drivers.DrawAnimation(greenHill, "flower2", animationIndex+7, 220, 115)
		drivers.DrawAnimation(greenHill, "flower2", animationIndex, 250, 115)

		RenderHUD(animationIndex, selectedVMIndex)

		drivers.DrawAnimation(sonic, "sonic", animationIndex, 35, 128)


	*/
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
			model.SwitchToVM(centinelVM, vms[selectedVMIndex])
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
