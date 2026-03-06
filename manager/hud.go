package manager

import (
	"time"

	"apodeiktikos.com/fbtest/core"
	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/model"
	"apodeiktikos.com/fbtest/util"
)

const hudOffset float64 = 175

var ring *engine.SpriteInstance
var gpuString string
var centinelVM *model.VM
var vms []model.VM
var selectedVMIndex = 0
var texts []*engine.Text

func SetupHud(sprites core.Sprites, renderables []engine.Renderable) []engine.Renderable {
	for y := 0; y < 13; y++ {
		renderables = append(renderables, engine.BuildSpriteInstance(sprites, "ZigZag", "idle", core.Point{X: hudOffset, Y: float64(y * 16)}))
	}

	texts = SetupHUDTexts(sprites, hudOffset)
	for _, t := range texts {
		renderables = append(renderables, t) // Aquí 't' se convierte a Renderable automáticamente
	}

	renderables = append(renderables, SetupRing(sprites)) // Aquí 't' se convierte a Renderable automáticamente

	return renderables
}

func SetupHUDTexts(sprites core.Sprites, hudOffset float64) []*engine.Text {
	texts := make([]*engine.Text, 0)
	gpuString = util.ContextStorage.GpuString
	centinelVM = model.GetVMByName(util.ContextStorage.CentineVMName)
	vms = model.GetVMsWithGPU(gpuString, centinelVM)
	for i, vm := range vms {
		var text string
		if vm.Equals(centinelVM) {
			text = "POWER OFF"
		} else {
			text = vm.Name
		}

		textInstance := engine.BuildTextInstance(sprites.Sprites["BoldLetters"], text, core.Point{X: hudOffset + 30, Y: float64(i)*16 + 60})
		texts = append(texts, textInstance)
	}

	// ToDo uncomment next lines after rebugging
	texts[0].Position.X += 12
	texts[len(texts)-1].Position.Y += 12

	texts = append(texts, engine.BuildTextInstance(sprites.Sprites["GenesisLetters"], centinelVM.Name, core.Point{X: hudOffset + 20, Y: 30}))

	return texts
}

func IncrementVMIndex(value int) {
	if value == 0 || selectedVMIndex >= len(vms) || selectedVMIndex < 0 {
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
	ring.OnMovementComplete = OnMovementComplete
}

func SelectMenuOption() {
	ci := 0
	for i, text := range texts[:len(texts)-1] {
		if i == selectedVMIndex || i == len(texts)-1 {
			continue
		}
		text.OnMovementComplete = OnMovementComplete
		text.SetEaseFunction(util.EaseInOutQuad)
		text.MoveTo(text.Position, time.Duration(ci*100)*time.Millisecond)
		ci++
	}
}

func OnRingAnimationComplete(sprite *engine.SpriteInstance) {
	if sprite == ring {
		fadeSeq := ring.Sprite.Sequences["fade"]
		if &ring.CurrentSequence[0] == &fadeSeq[0] {
			ring.CurrentSequence = ring.Sprite.Sequences["end"]
			model.SwitchToVM(centinelVM, vms[selectedVMIndex])
		}

	}
}

func OnMovementComplete(sprite engine.Renderable) {
	for _, text := range texts {
		if sprite == text {
			text.SetEaseFunction(util.EaseInOutQuad)
			text.MoveTo(core.Point{X: 320, Y: text.Position.Y}, 300*time.Millisecond)
		}
	}
}
