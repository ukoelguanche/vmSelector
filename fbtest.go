package main

import (
	_ "image/png"
	"math"
	"time"

	"apodeiktikos.com/fbtest/drivers"
	"apodeiktikos.com/fbtest/loaders"
	"apodeiktikos.com/fbtest/model"
	"apodeiktikos.com/fbtest/util"
)

const targetFPS = 30
const frameDelay = time.Second / targetFPS

var gpuString string
var centinelVM *model.VM
var vms []model.VM

var sprites model.Sprites
var ring *model.SpriteInstance
var spriteInstances []*model.SpriteInstance
var texts []*model.Text

var selectedVMIndex = 0

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

	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Sonic", "idle", model.Point{X: 35, Y: 131}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillForeground", "idle", model.Point{X: 0, Y: 0}))

	const hudOffset int32 = 155

	for y := 0; y < 13; y++ {
		spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "ZigZag", "idle", model.Point{X: hudOffset, Y: int32(y * 16)}))
	}
	ring = model.BuildSpriteInstance(sprites, "Ring", "idle", model.Point{X: hudOffset + 30, Y: 56})
	spriteInstances = append(spriteInstances, ring)

	texts = make([]*model.Text, 0)

	gpuString = util.ContextStorage.GpuString
	centinelVM = model.GetVMByName(util.ContextStorage.CentineVMName)
	vms = model.GetVMsWithGPU(gpuString, centinelVM)
	for i, vm := range vms {
		var text string
		if vm.Equals(centinelVM) {
			text = "Power off"
		} else {
			text = vm.Name
		}

		textInstance := model.BuildTextInstance(sprites.Sprites["BoldLetters"], text, model.Point{X: hudOffset + 36, Y: int32(i)*16 + 60})
		texts = append(texts, textInstance)
	}

	texts[0].Position.X += 12
	texts[0].TargetPosition.X += 12
	texts = append(texts, model.BuildTextInstance(sprites.Sprites["GenesisLetters"], centinelVM.Name, model.Point{X: hudOffset + 30, Y: 30}))

	return

}

func Loop(animationIndex int, selectedVMIndex int, endLoop bool) {
	drivers.GlobalDisplay.Clear()

	for _, spriteInstance := range spriteInstances {
		drivers.DrawAnimation(spriteInstance)
		spriteInstance.NextFrame()
	}

	for _, text := range texts {
		drivers.DrawText(text)
		text.NextFrame()
	}

	drivers.GlobalDisplay.Present()

}

func incrementVMIndex(value int) {
	if math.Abs(float64(ring.TargetPosition.Y-ring.Position.Y)) > 1 {
		return
	}
	texts[selectedVMIndex].TargetPosition.X -= 14
	selectedVMIndex = max(0, min(len(vms)-1, selectedVMIndex+value))
	texts[selectedVMIndex].TargetPosition.X += 14

	ring.TargetPosition.Y = texts[selectedVMIndex].Position.Y - 4
	ring.Speed = 1
}

func main() {
	Init()

	defer drivers.GlobalDisplay.Close()
	var animationIndex = 0

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
			incrementVMIndex(1)
		}
		if (dx < 0 || dy < 0) && selectedVMIndex > 0 {
			incrementVMIndex(-1)
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
