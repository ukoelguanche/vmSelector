package main

import (
	"fmt"
	_ "image/png"
	"math"
	"math/rand"
	"time"

	"apodeiktikos.com/fbtest/drivers"
	"apodeiktikos.com/fbtest/loaders"
	"apodeiktikos.com/fbtest/model"
	"apodeiktikos.com/fbtest/util"
)

const targetFPS = 25
const frameDelay = time.Second / targetFPS

var gpuString string
var centinelVM *model.VM
var vms []model.VM

var sprites model.Sprites
var ring *model.SpriteInstance
var sonic *model.SpriteInstance
var clouds []*model.SpriteInstance

var spriteInstances []*model.SpriteInstance
var texts []*model.Text

var selectedVMIndex = 0

const hudOffset float64 = 175

func Init() {
	util.LoadContext()
	drivers.GlobalDisplay = drivers.InitDisplay(drivers.VW, drivers.VH)
	drivers.GlobalKeyboard = drivers.InitKeyboard()

	loaders.LoadSprites("./resources/sprites/Sprites.json", &sprites)

	spriteInstances = make([]*model.SpriteInstance, 0)

	SetupClouds()
	SetupGreenHillBackground()
	SetupSonic()
	SetupGreenHillForeground()
	SetupHud()

	spriteInstances = append(spriteInstances, ring)
}

func SetupSonic() {
	sonic = model.BuildSpriteInstance(sprites, "Sonic", "idle", model.Point{X: 35, Y: 131})
	sonic.OnAnimationComplete = OnAnimationComplete
	spriteInstances = append(spriteInstances, sonic)
}

func SetupGreenHillForeground() {
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower1", "idle", model.Point{X: 154, Y: 90}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower2", "idle", model.Point{X: -5, Y: 115}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower2", "idle", model.Point{X: 220, Y: 115}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower2", "idle", model.Point{X: 250, Y: 115}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillForeground", "idle", model.Point{X: 0, Y: 0}))
}

func SetupHud() {
	for y := 0; y < 13; y++ {
		spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "ZigZag", "idle", model.Point{X: hudOffset, Y: float64(y * 16)}))
	}

	SetupHUDTexts(hudOffset)

	ring = model.BuildSpriteInstance(sprites, "Ring", "idle", model.Point{X: hudOffset + 20, Y: 56})
	ring.SetEaseFunction(util.EaseInOutCubic)
	ring.OnAnimationComplete = OnAnimationComplete
}

func SetupHUDTexts(hudOffset float64) {
	texts = make([]*model.Text, 0)
	gpuString = util.ContextStorage.GpuString
	centinelVM = model.GetVMByName(util.ContextStorage.CentineVMName)
	vms = model.GetVMsWithGPU(gpuString, centinelVM)
	texts = append(texts, model.BuildTextInstance(sprites.Sprites["GenesisLetters"], centinelVM.Name, model.Point{X: hudOffset + 20, Y: 30}))
	for i, vm := range vms {
		var text string
		if vm.Equals(centinelVM) {
			text = "POWER OFF"
		} else {
			text = vm.Name
		}

		textInstance := model.BuildTextInstance(sprites.Sprites["BoldLetters"], text, model.Point{X: hudOffset + 30, Y: float64(i)*16 + 60})
		texts = append(texts, textInstance)
	}

	texts[0].Position.X += 12
	texts[len(texts)-1].Position.Y += 12
}

func SetupGreenHillBackground() {
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer4", "idle", model.Point{X: 0, Y: 64}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer5", "idle", model.Point{X: 0, Y: 112}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer6", "idle", model.Point{X: 0, Y: 152}))
}

func SetupClouds() {
	var cloudSprite *model.SpriteInstance
	clouds = make([]*model.SpriteInstance, 0)
	var y float64 = 0
	for i := 0; i < 3; i++ {
		cloudSprite = model.BuildSpriteInstance(sprites, fmt.Sprintf("GreenHillBackgroundLayer%d", i+1), "idle", model.Point{X: 0, Y: float64(y)})

		cloudSprite.MoveTo(cloudSprite.Position.SetX(-980), time.Duration(90000+i*10000)*time.Millisecond)

		cloudSprite.OnMovementComplete = OnMovementComplete
		spriteInstances = append(spriteInstances, cloudSprite)
		y += cloudSprite.Sprite.Frames[0].Size.H

		clouds = append(clouds, cloudSprite)
	}
}

func Loop() {
	drivers.GlobalDisplay.Clear()

	for _, spriteInstance := range spriteInstances {
		drivers.DrawSpriteFrame(spriteInstance)
		spriteInstance.NextFrame()
	}

	for _, text := range texts {
		drivers.DrawText(text)
		text.NextFrame()
	}

	drivers.GlobalDisplay.Present()
}

func OnAnimationComplete(sprite *model.SpriteInstance) {
	if sprite == sonic {
		keys := make([]string, 0, len(sprite.Sprite.Sequences))
		for k := range sprite.Sprite.Sequences {
			keys = append(keys, k)
		}
		randomSequence := keys[rand.Intn(len(keys))]
		sonic.CurrentSequence = sonic.Sprite.Sequences[randomSequence]
	}

	if sprite == ring {
		fadeSeq := ring.Sprite.Sequences["fade"]
		if &ring.CurrentSequence[0] == &fadeSeq[0] {
			ring.CurrentSequence = ring.Sprite.Sequences["end"]
			model.SwitchToVM(centinelVM, vms[selectedVMIndex])
		}

	}
}

func OnMovementComplete(sprite *model.SpriteInstance) {
	if sprite == clouds[0] || sprite == clouds[1] || sprite == clouds[2] {
		sprite.Position = sprite.Position.SetX(0)
		sprite.SetTargetPosition(sprite.Position.SetX(-980))
	}
}

func incrementVMIndex(value int) {
	if value == 0 || selectedVMIndex >= len(vms) || selectedVMIndex < 0 {
		return
	}
	if math.Abs(float64(ring.TargetPosition.Y-ring.Position.Y)) > 1 {
		return
	}
	if ring.IsMoving() {
		return
	}

	const transitionDuaration = 200 * time.Millisecond
	texts[selectedVMIndex].MoveTo(texts[selectedVMIndex].Position.SetX(hudOffset+30), transitionDuaration)
	selectedVMIndex = max(0, min(len(vms)-1, selectedVMIndex+value))
	texts[selectedVMIndex].MoveTo(texts[selectedVMIndex].Position.SetX(hudOffset+42), transitionDuaration)

	ring.MoveTo(ring.Position.SetY(texts[selectedVMIndex].Position.Y-4), transitionDuaration)
}

func handleKeyboardInput() bool {
	kbd := drivers.GlobalKeyboard.GetInput()

	if kbd == drivers.KBD_ESCAPE {
		return true
	}

	if kbd == drivers.KBD_RETURN {
		ring.CurrentSequencePosition = 0.0
		for i, text := range texts {
			if i == selectedVMIndex || i == 0 {
				continue
			}
			text.SetEaseFunction(util.EaseInOutCubic)
			text.MoveTo(model.Point{X: 350, Y: text.Position.Y}, time.Duration((i+1)*300)*time.Millisecond)
		}
	}

	var inc int
	if kbd == drivers.KBD_UP || kbd == drivers.KBD_LEFT {
		inc = -1
	} else if kbd == drivers.KBD_DOWN || kbd == drivers.KBD_RIGHT {
		inc = 1
	}

	incrementVMIndex(inc)

	return false
}

func main() {
	Init()

	defer drivers.GlobalDisplay.Close()
	defer drivers.GlobalDisplay.Clear()
	defer drivers.GlobalKeyboard.Close()

	for {
		start := time.Now()

		if handleKeyboardInput() {
			break
		}

		Loop()

		elapsed := time.Since(start)
		if elapsed < frameDelay {
			time.Sleep(frameDelay - elapsed)
		}
	}
}
