package main

import (
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
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
var sonic *model.SpriteInstance

var spriteInstances []*model.SpriteInstance
var texts []*model.Text

var selectedVMIndex = 0

func Init() {
	util.LoadContext()
	drivers.GlobalDisplay = drivers.InitDisplay(drivers.VW, drivers.VH)

	loaders.LoadSprites("./resources/sprites/Sprites.json", &sprites)

	spriteInstances = make([]*model.SpriteInstance, 0)

	layer1 := model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer1", "idle", model.Point{X: 0, Y: 0})
	layer1.TargetPosition.X = -3000
	layer1.Speed = 3
	spriteInstances = append(spriteInstances, layer1)

	layer2 := model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer2", "idle", model.Point{X: 0, Y: 32})
	layer2.TargetPosition.X = -3000
	layer2.Speed = 2
	spriteInstances = append(spriteInstances, layer2)
	layer3 := model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer3", "idle", model.Point{X: 0, Y: 48})
	layer3.TargetPosition.X = -3000
	layer3.Speed = 1
	spriteInstances = append(spriteInstances, layer3)
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer4", "idle", model.Point{X: 0, Y: 64}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer5", "idle", model.Point{X: 0, Y: 112}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer6", "idle", model.Point{X: 0, Y: 152}))

	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower1", "idle", model.Point{X: 154, Y: 90}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower2", "idle", model.Point{X: -5, Y: 115}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower2", "idle", model.Point{X: 220, Y: 115}))
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "Flower2", "idle", model.Point{X: 250, Y: 115}))

	sonic = model.BuildSpriteInstance(sprites, "Sonic", "idle", model.Point{X: 35, Y: 131})
	sonic.OnComplete = OnComplete
	spriteInstances = append(spriteInstances, sonic)
	spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "GreenHillForeground", "idle", model.Point{X: 0, Y: 0}))

	const hudOffset int32 = 175

	for y := 0; y < 13; y++ {
		spriteInstances = append(spriteInstances, model.BuildSpriteInstance(sprites, "ZigZag", "idle", model.Point{X: hudOffset, Y: int32(y * 16)}))
	}

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

		textInstance := model.BuildTextInstance(sprites.Sprites["BoldLetters"], text, model.Point{X: hudOffset + 30, Y: int32(i)*16 + 60})
		texts = append(texts, textInstance)
	}

	texts[0].Position.X += 12
	texts[0].TargetPosition.X += 12
	texts = append(texts, model.BuildTextInstance(sprites.Sprites["GenesisLetters"], centinelVM.Name, model.Point{X: hudOffset + 20, Y: 30}))

	ring = model.BuildSpriteInstance(sprites, "Ring", "idle", model.Point{X: hudOffset + 20, Y: 56})
	ring.OnComplete = OnComplete
	spriteInstances = append(spriteInstances, ring)

	return

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

func OnComplete(sprite *model.SpriteInstance) {
	if sprite == sonic {
		keys := make([]string, 0, len(sprite.Sprite.Sequences))
		for k := range sprite.Sprite.Sequences {
			keys = append(keys, k)
		}
		randomSequence := keys[rand.Intn(len(keys))]
		sonic.CurrentSequence = sonic.Sprite.Sequences[randomSequence]
	} else if sprite == ring {
		fadeSeq := ring.Sprite.Sequences["fade"]
		if &ring.CurrentSequence[0] == &fadeSeq[0] {
			ring.CurrentSequence = ring.Sprite.Sequences["end"]
			model.SwitchToVM(centinelVM, vms[selectedVMIndex])
		}

	}
}

func incrementVMIndex(value int) {
	if selectedVMIndex >= len(vms) || selectedVMIndex < 0 {
		return
	}

	if math.Abs(float64(ring.TargetPosition.Y-ring.Position.Y)) > 1 {
		return
	}
	texts[selectedVMIndex].TargetPosition.X -= 12
	selectedVMIndex = max(0, min(len(vms)-1, selectedVMIndex+value))
	texts[selectedVMIndex].TargetPosition.X += 12

	ring.TargetPosition.Y = texts[selectedVMIndex].Position.Y - 4
	ring.Speed = 2
}

func handleKeyboardInput() bool {
	dx, quit, enter := drivers.GlobalDisplay.GetInput()

	if quit {
		return true
	}

	if enter {
		ring.CurrentSequencePosition = 0.0
		ring.CurrentSequence = ring.Sprite.Sequences["fade"]
	}

	incrementVMIndex(int(dx))

	return false
}

func waitForExit(cleanup func()) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan,
		os.Interrupt,    // Ctrl+C
		syscall.SIGTERM, // kill
	)

	<-sigChan
	cleanup()
	os.Exit(0)
}

func main() {
	Init()

	go waitForExit(func() {
		drivers.GlobalDisplay.Close()
		drivers.GlobalDisplay.Clear()
	})

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
